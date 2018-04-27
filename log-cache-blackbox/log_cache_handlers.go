package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	logcache "code.cloudfoundry.org/go-log-cache"
	"code.cloudfoundry.org/go-loggregator/rpc/loggregator_v2"
)

func latencyHandler(cfg Config, httpClient *http.Client) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithCancel(r.Context())
		defer cancel()

		client := logcache.NewClient(
			cfg.LogCacheURL.String(),
			logcache.WithHTTPClient(
				logcache.NewOauth2HTTPClient(
					cfg.UAAAddr,
					cfg.UAAClient,
					cfg.UAAClientSecret,
					logcache.WithOauth2HTTPClient(httpClient),
				),
			),
		)

		latCtx, _ := context.WithTimeout(ctx, 10*time.Second)
		resultData, err := measureLatency(latCtx, client.Read, cfg.VCapApp.ApplicationID)
		if err != nil {
			log.Printf("error getting result data: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Write(resultData)
		w.Write([]byte("\n"))
	})
}

func reliabilityHandler(cfg Config, httpClient *http.Client) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithCancel(r.Context())
		defer cancel()
		client := buildClient(cfg, httpClient)
		reader := client.Read

		primeCtx, _ := context.WithTimeout(ctx, time.Minute)
		err := prime(primeCtx, cfg.VCapApp.ApplicationID, reader)
		if err != nil {
			log.Printf("unable to prime for source id: %s: %s", cfg.VCapApp.ApplicationID, err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		emitCount := 10000
		prefix := fmt.Sprintf("%d - ", time.Now().UnixNano())

		start := time.Now()
		for i := 0; i < emitCount; i++ {
			log.Printf("%s%d", prefix, i)
			time.Sleep(time.Millisecond)
		}
		end := time.Now()

		var (
			receivedCount    int
			badReceivedCount int
		)
		walkCtx, _ := context.WithTimeout(ctx, 40*time.Second)
		logcache.Walk(
			walkCtx,
			cfg.VCapApp.ApplicationID,
			func(envelopes []*loggregator_v2.Envelope) bool {
				for _, e := range envelopes {
					if strings.Contains(string(e.GetLog().GetPayload()), prefix) {
						receivedCount++
					} else {
						badReceivedCount++
					}
				}
				return receivedCount < emitCount
			},
			reader,
			logcache.WithWalkStartTime(start),
			logcache.WithWalkEndTime(end),
			logcache.WithWalkBackoff(logcache.NewAlwaysRetryBackoff(time.Second)),
		)

		result := ReliabilityTestResult{
			LogsSent:        emitCount,
			LogsReceived:    receivedCount,
			BadLogsReceived: badReceivedCount,
		}

		resultData, err := json.Marshal(&result)
		if err != nil {
			log.Printf("failed to marshal test results: %s", err)
			return
		}

		w.Write(resultData)
		w.Write([]byte("\n"))
	})
}

func buildClient(cfg Config, httpClient *http.Client) *logcache.Client {
	return logcache.NewClient(
		cfg.LogCacheURL.String(),
		logcache.WithHTTPClient(logcache.NewOauth2HTTPClient(
			cfg.UAAAddr,
			cfg.UAAClient,
			cfg.UAAClientSecret,
			logcache.WithOauth2HTTPClient(httpClient)),
		),
	)
}

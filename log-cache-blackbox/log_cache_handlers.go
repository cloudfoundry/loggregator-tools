package main

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	logcache "code.cloudfoundry.org/go-log-cache"
	"code.cloudfoundry.org/go-loggregator/rpc/loggregator_v2"
)

func reliabilityHandler(cfg Config) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		emitCount := 10000
		prefix := fmt.Sprintf("%d - ", time.Now().UnixNano())

		start := time.Now()
		for i := 0; i < emitCount; i++ {
			log.Printf("%s %d", prefix, i)
			time.Sleep(time.Millisecond)
		}
		end := time.Now()

		// Give the system time to get the envelopes
		time.Sleep(10 * time.Second)

		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: cfg.SkipSSLValidation},
		}

		client := logcache.NewClient(cfg.LogCacheAddr,
			logcache.WithHTTPClient(logcache.NewOauth2HTTPClient(
				cfg.UAAAddr,
				cfg.UAAClient,
				cfg.UAAClientSecret,
				logcache.WithOauth2HTTPClient(
					&http.Client{
						Transport: tr,
						Timeout:   5 * time.Second,
					},
				))),
		)

		var receivedCount int
		logcache.Walk(context.Background(), cfg.VCapApp.ApplicationID, func(envelopes []*loggregator_v2.Envelope) bool {
			for _, e := range envelopes {
				if strings.Contains(string(e.GetLog().GetPayload()), prefix) {
					receivedCount++
				}
			}
			return receivedCount < emitCount
		},
			client.Read,
			logcache.WithWalkStartTime(start),
			logcache.WithWalkEndTime(end),
			logcache.WithWalkBackoff(logcache.NewRetryBackoff(50*time.Millisecond, 100)),
		)

		result := ReliabilityTestResult{
			LogsSent:     emitCount,
			LogsReceived: receivedCount,
		}

		resultData, err := json.Marshal(&result)
		if err != nil {
			log.Printf("failed to marshal test results: %s", err)
			return
		}

		w.Write(resultData)
	})
}

func latencyHandler(cfg Config) http.Handler {
	httpClient := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: cfg.SkipSSLValidation},
		},
		Timeout: 5 * time.Second,
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		client := logcache.NewClient(cfg.LogCacheAddr,
			logcache.WithHTTPClient(logcache.NewOauth2HTTPClient(
				cfg.UAAAddr,
				cfg.UAAClient,
				cfg.UAAClientSecret,
				logcache.WithOauth2HTTPClient(httpClient),
			)),
		)

		resultData, err := measureLatency(client.Read, cfg.VCapApp.ApplicationID)
		if err != nil {
			log.Printf("error getting result data: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Write(resultData)
		return
	})
}

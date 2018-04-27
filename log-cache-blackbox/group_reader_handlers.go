package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"

	logcache "code.cloudfoundry.org/go-log-cache"
	"code.cloudfoundry.org/go-loggregator/rpc/loggregator_v2"
	uuid "github.com/nu7hatch/gouuid"
)

func groupLatencyHandler(cfg Config, httpClient *http.Client) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithCancel(r.Context())
		defer cancel()
		client := buildGroupClient(cfg, httpClient)

		groupUUID, err := uuid.NewV4()
		if err != nil {
			log.Printf("unable to create groupUUID: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		groupName := groupUUID.String()

		latCtx, _ := context.WithTimeout(ctx, 10*time.Second)
		resultData, err := measureLatency(
			latCtx,
			client.BuildReader(rand.Uint64()),
			groupName,
		)
		if err != nil {
			log.Printf("error getting result data: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Write(resultData)
		w.Write([]byte("\n"))
	})
}

func groupReliabilityHandler(cfg Config, httpClient *http.Client) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithCancel(r.Context())
		defer cancel()
		client := buildGroupClient(cfg, httpClient)
		reader := client.BuildReader(rand.Uint64())

		size, err := strconv.Atoi(r.URL.Query().Get("size"))
		if err != nil {
			log.Printf("invalid size: %s %s", r.URL.Query().Get("size"), err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		groupUUID, err := uuid.NewV4()
		if err != nil {
			log.Printf("unable to create groupUUID: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		groupName := groupUUID.String()

		sIDs, err := sourceIDs(httpClient, cfg, size)
		if err != nil {
			log.Printf("unable to get sourceIDs: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		go maintainGroup(ctx, groupName, sIDs, client)

		primeCtx, _ := context.WithTimeout(ctx, time.Minute)
		err = prime(primeCtx, groupName, reader)
		if err != nil {
			log.Printf("unable to prime for group: %s: %s", groupName, err)
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
			groupName,
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

func buildGroupClient(cfg Config, httpClient *http.Client) *logcache.ShardGroupReaderClient {
	return logcache.NewShardGroupReaderClient(
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
}

func maintainGroup(
	ctx context.Context,
	groupName string,
	sIDs []string,
	client *logcache.ShardGroupReaderClient,
) {
	ticker := time.NewTicker(10 * time.Second)
	for {
		for _, sID := range sIDs {
			shardGroupCtx, _ := context.WithTimeout(ctx, time.Second)
			err := client.SetShardGroup(shardGroupCtx, groupName, sID)
			if err != nil {
				log.Printf("unable to set shard group: %s", err)
			}
		}

		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			continue
		}
	}
}

func sourceIDs(httpClient *http.Client, cfg Config, size int) ([]string, error) {
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
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	meta, err := client.Meta(ctx)
	if err != nil {
		return nil, err
	}

	sourceIDs := make([]string, 0, size)
	for k := range meta {
		if k == cfg.VCapApp.ApplicationID {
			continue
		}
		if len(sourceIDs) < size-1 {
			sourceIDs = append(sourceIDs, k)
		}
	}
	sourceIDs = append(sourceIDs, cfg.VCapApp.ApplicationID)
	if len(sourceIDs) != size {
		return nil, fmt.Errorf("Asked for %d source IDs but only found %d", size, len(sourceIDs))
	}
	return sourceIDs, nil
}

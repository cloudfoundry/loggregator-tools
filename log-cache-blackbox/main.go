package main

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync/atomic"
	"time"

	envstruct "code.cloudfoundry.org/go-envstruct"
	logcache "code.cloudfoundry.org/go-log-cache"
	"code.cloudfoundry.org/go-loggregator/rpc/loggregator_v2"
	uuid "github.com/nu7hatch/gouuid"
)

type Config struct {
	LogCacheAddr string `env:"LOG_CACHE_ADDR,   required"`

	Port    int     `env:"PORT,             required"`
	VCapApp VCapApp `env:"VCAP_APPLICATION, required"`

	UAAAddr         string `env:"UAA_ADDR,          required"`
	UAAClient       string `env:"UAA_CLIENT,        required"`
	UAAClientSecret string `env:"UAA_CLIENT_SECRET, required, noreport"`

	SkipSSLValidation bool `env:"SKIP_SSL_VALIDATION"`
}

type VCapApp struct {
	ApplicationID string `json:"application_id"`
}

func (v *VCapApp) UnmarshalEnv(jsonData string) error {
	return json.Unmarshal([]byte(jsonData), &v)
}

type LatencyTestResult struct {
	Latency          float64   `json:"latency"`
	QueryTimes       []float64 `json:"queryTime"`
	AverageQueryTime float64   `json:"averageQueryTime"`
}

type ReliabilityTestResult struct {
	LogsSent     int `json:"logsSent"`
	LogsReceived int `json:"logsReceived"`
}

func main() {
	var cfg Config
	err := envstruct.Load(&cfg)
	if err != nil {
		log.Fatalf("failed to load config: %s", err)
	}

	err = envstruct.WriteReport(&cfg)
	if err != nil {
		log.Fatalf("failed to write config report: %s", err)
	}

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", cfg.Port), handler(cfg)))
}

func handler(cfg Config) http.Handler {
	var testRunning int64

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		if !atomic.CompareAndSwapInt64(&testRunning, 0, 1) {
			w.WriteHeader(http.StatusLocked)
			return
		}
		defer atomic.StoreInt64(&testRunning, 0)

		switch r.URL.Path {
		case "/":
			latencyHandler(cfg).ServeHTTP(w, r)
		case "/group":
			groupLatencyHandler(cfg).ServeHTTP(w, r)
		case "/reliability":
			reliabilityHandler(cfg).ServeHTTP(w, r)
		case "/group/reliability":
			groupReliabilityHandler(cfg).ServeHTTP(w, r)
		default:
			w.WriteHeader(http.StatusNotFound)
			return
		}
	})
}

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

func groupReliabilityHandler(cfg Config) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		reader, err := buildGroupReader(ctx, size, groupName, cfg)
		if err != nil {
			log.Printf("Unable to create group reader: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		emitCount := 10000
		prefix := fmt.Sprintf("%d - ", time.Now().UnixNano())

		start := time.Now()
		go func() {
			// Give the system time to get the envelopes
			time.Sleep(20 * time.Second)

			for i := 0; i < emitCount; i++ {
				log.Printf("%s %d", prefix, i)
				time.Sleep(time.Millisecond)
			}
		}()

		var receivedCount int
		walkCtx, _ := context.WithTimeout(ctx, 40*time.Second)
		log.Printf("Starting walk...")
		logcache.Walk(
			walkCtx,
			groupName,
			func(envelopes []*loggregator_v2.Envelope) bool {
				for _, e := range envelopes {
					if strings.Contains(string(e.GetLog().GetPayload()), prefix) {
						receivedCount++
					}
				}
				return receivedCount < emitCount
			},
			reader,
			logcache.WithWalkStartTime(start),
			logcache.WithWalkBackoff(logcache.NewRetryBackoff(50*time.Millisecond, 100)),
			logcache.WithWalkLogger(log.New(os.Stderr, "[WALK]", 0)),
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

func buildGroupReader(ctx context.Context, size int, groupName string, cfg Config) (logcache.Reader, error) {
	httpClient := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: cfg.SkipSSLValidation},
		},
		Timeout: 5 * time.Second,
	}

	sIDs, err := sourceIDs(httpClient, cfg, size)
	if err != nil {
		return nil, fmt.Errorf("unable to get sourceIDs: %s", err)
	}

	client := logcache.NewShardGroupReaderClient(
		cfg.LogCacheAddr,
		logcache.WithHTTPClient(
			logcache.NewOauth2HTTPClient(
				cfg.UAAAddr,
				cfg.UAAClient,
				cfg.UAAClientSecret,
				logcache.WithOauth2HTTPClient(httpClient),
			),
		),
	)

	for _, sID := range sIDs {
		go func(sID string) {
			ticker := time.NewTicker(time.Second)
			for {
				ctx, _ := context.WithTimeout(ctx, 10*time.Second)
				err = client.SetShardGroup(ctx, groupName, sID)
				if err != nil {
					log.Printf("unable to set shard group: %s", err)
				}

				select {
				case <-ctx.Done():
					return
				case <-ticker.C:
					continue
				}
			}
		}(sID)
	}

	return client.BuildReader(rand.Uint64()), nil
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

func groupLatencyHandler(cfg Config) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		reader, err := buildGroupReader(ctx, size, groupName, cfg)
		if err != nil {
			log.Printf("Unable to create group reader: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		resultData, err := measureLatency(reader, groupName)
		if err != nil {
			log.Printf("error getting result data: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Write(resultData)
		return
	})
}

func measureLatency(reader logcache.Reader, name string) ([]byte, error) {
	var queryTimes []time.Duration
	for j := 0; j < 10; j++ {
		expectedLog := fmt.Sprintf("Test log - %d", rand.Int63())
		logStartTime := time.Now()
		fmt.Println(expectedLog)

		for i := 0; i < 100; i++ {

			// Look at 110 pages. This equates to 11000 possible envelopes. The
			// reliability test emits 10000 envelopes and could be clutter within the
			// cache.
			queryDurations, success := consumePages(
				110,
				logStartTime.Add(-7500*time.Millisecond),
				logStartTime.Add(7500*time.Millisecond),
				reader,
				name,
				func(envelopes []*loggregator_v2.Envelope) (bool, int64) {
					return lookForMsg(expectedLog, envelopes)
				},
			)
			queryTimes = append(queryTimes, queryDurations...)

			if success {
				var queryTimesMs []float64
				var totalQueryTimes time.Duration
				for _, qt := range queryTimes {
					totalQueryTimes += qt
					queryTimesMs = append(queryTimesMs, float64(qt)/float64(time.Millisecond))
				}

				avgQT := int64(totalQueryTimes) / int64(len(queryTimes))

				// Success
				testResults := LatencyTestResult{
					Latency:          float64(time.Since(logStartTime)) / float64(time.Millisecond),
					QueryTimes:       queryTimesMs,
					AverageQueryTime: float64(avgQT) / float64(time.Millisecond),
				}

				resultData, err := json.Marshal(testResults)
				if err != nil {
					return nil, err
				}

				return resultData, nil
			}
		}
	}

	// Fail...
	return nil, errors.New("Never found the log")
}

func consumePages(pages int, start, end time.Time, reader logcache.Reader, name string, f func([]*loggregator_v2.Envelope) (bool, int64)) ([]time.Duration, bool) {
	var queryDurations []time.Duration

	for i := 0; i < pages; i++ {
		queryDuration, last, count, err := consumeTimeRange(start, end, reader, name, f)
		if err != nil {
			log.Println(err)
			return nil, false
		}

		queryDurations = append(queryDurations, queryDuration)
		if !last.IsZero() && count > 0 {
			start = last.Add(time.Nanosecond)
			continue
		}

		return queryDurations, true
	}

	return queryDurations, false
}

func consumeTimeRange(start, end time.Time, reader logcache.Reader, name string, f func([]*loggregator_v2.Envelope) (bool, int64)) (time.Duration, time.Time, int, error) {
	queryStart := time.Now()
	data, err := reader(
		context.Background(),
		name,
		start,
		logcache.WithEndTime(end),
	)

	if err != nil {
		return 0, time.Time{}, 0, err
	}
	queryDuration := time.Since(queryStart)

	found, last := f(data)
	if !found {
		return queryDuration, time.Unix(0, last), len(data), nil
	}

	return queryDuration, time.Time{}, len(data), nil
}

func lookForMsg(msg string, envelopes []*loggregator_v2.Envelope) (bool, int64) {
	var last int64
	for _, e := range envelopes {
		last = e.Timestamp
		if string(e.GetLog().GetPayload()) == msg {
			return true, last
		}
	}

	return false, last
}

func sourceIDs(httpClient *http.Client, cfg Config, size int) ([]string, error) {
	client := logcache.NewClient(
		cfg.LogCacheAddr,
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
	return sourceIDs, nil
}

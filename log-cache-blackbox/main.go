package main

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"sync/atomic"
	"time"

	envstruct "code.cloudfoundry.org/go-envstruct"
	logcache "code.cloudfoundry.org/go-log-cache"
	"code.cloudfoundry.org/go-loggregator/rpc/loggregator_v2"
)

type Config struct {
	LogCacheAddr string `env:"LOG_CACHE_ADDR,   required"`

	Port    int     `env:"PORT,             required"`
	VCapApp VCapApp `env:"VCAP_APPLICATION, required"`

	UAAAddr         string `env:"UAA_ADDR,          required"`
	UAAClient       string `env:"UAA_CLIENT,        required"`
	UAAClientSecret string `env:"UAA_CLIENT_SECRET, required"`

	SkipSSLValidation bool `env:"SKIP_SSL_VALIDATION"`
}

type VCapApp struct {
	ApplicationID string `json:"application_id"`
}

func (v *VCapApp) UnmarshalEnv(jsonData string) error {
	return json.Unmarshal([]byte(jsonData), &v)
}

type TestResult struct {
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
		case "/reliability":
			reliabilityHandler(cfg).ServeHTTP(w, r)
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
			receivedCount += len(envelopes)
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

func waitForEnvelopes(ctx context.Context, cfg Config, emitCount int, prefix string) int {
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
	for {
		select {
		case <-ctx.Done():
			return receivedCount
		default:
			data, err := client.Read(context.Background(), cfg.VCapApp.ApplicationID, time.Unix(0, 0))
			if err != nil {
				log.Printf("error while reading: %s", err)
				continue
			}

			receivedCount = countMessages(prefix, data)
			if receivedCount == emitCount {
				return receivedCount
			}
		}
	}

	return receivedCount
}

func latencyHandler(cfg Config) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
					cfg,
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
					testResults := TestResult{
						Latency:          float64(time.Since(logStartTime)) / float64(time.Millisecond),
						QueryTimes:       queryTimesMs,
						AverageQueryTime: float64(avgQT) / float64(time.Millisecond),
					}

					resultData, err := json.Marshal(testResults)
					if err != nil {
						log.Printf("failed to marshal test results: %s", err)
						w.WriteHeader(http.StatusInternalServerError)
						return
					}

					w.Write(resultData)
					return
				}
			}
		}

		// Fail...
		log.Println("Never found the log")
		w.WriteHeader(http.StatusInternalServerError)
	})
}

func consumePages(pages int, start, end time.Time, cfg Config, f func([]*loggregator_v2.Envelope) (bool, int64)) ([]time.Duration, bool) {
	var queryDurations []time.Duration

	for i := 0; i < pages; i++ {
		queryDuration, last, count, err := consumeTimeRange(start, end, cfg, f)
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

func consumeTimeRange(start, end time.Time, cfg Config, f func([]*loggregator_v2.Envelope) (bool, int64)) (time.Duration, time.Time, int, error) {
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

	queryStart := time.Now()
	data, err := client.Read(
		context.Background(),
		cfg.VCapApp.ApplicationID,
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

func countMessages(prefix string, envelopes []*loggregator_v2.Envelope) int {
	var count int
	for _, e := range envelopes {
		if strings.Contains(string(e.GetLog().GetPayload()), prefix) {
			count++
		}
	}

	return count
}

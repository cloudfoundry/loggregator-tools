package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"sync/atomic"
	"time"

	envstruct "code.cloudfoundry.org/go-envstruct"
)

type Config struct {
	Port         int     `env:"PORT,             required"`
	LogCacheAddr string  `env:"LOG_CACHE_ADDR,   required"`
	VCapApp      VCapApp `env:"VCAP_APPLICATION, required"`
}

type VCapApp struct {
	ApplicationID string `json:"application_id"`
}

func (v *VCapApp) UnmarshalEnv(jsonData string) error {
	return json.Unmarshal([]byte(jsonData), &v)
}

type LogCacheData struct {
	Envelopes []Envelope `json:"envelopes"`
}

type Envelope struct {
	Timestamp string `json:"timestamp"`
	Log       Log    `json:"log"`
}

type Log struct {
	Payload string `json:"payload"`
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
		log.Fatalf("failed to load config: %s", err)
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

		var receivedCount int
		consumePages(emitCount*100, start, end, cfg, func(envelopes []Envelope) (bool, int64) {
			if len(envelopes) == 0 {
				return false, 0
			}

			lastTs, err := strconv.ParseInt(envelopes[len(envelopes)-1].Timestamp, 10, 64)
			if err != nil {
				log.Printf("Failed to parse timestamp: %s", err)
				return false, 0
			}

			receivedCount += len(envelopes)
			return receivedCount >= emitCount, lastTs
		})

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
	url := fmt.Sprintf("%s/%s",
		cfg.LogCacheAddr,
		cfg.VCapApp.ApplicationID,
	)

	var receivedCount int
	for {
		select {
		case <-ctx.Done():
			return receivedCount
		default:
			resp, err := http.Get(url)
			if err != nil {
				log.Printf("error from log-cache: %s", err)
				continue
			}

			if resp.StatusCode != http.StatusOK {
				log.Printf("unxpected status code from log-cache: %d", resp.StatusCode)
				continue
			}

			var data LogCacheData
			err = json.NewDecoder(resp.Body).Decode(&data)
			if err != nil {
				log.Printf("failed to decode response body: %s", err)
				continue
			}

			receivedCount = countMessages(prefix, data.Envelopes)
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
					func(envelopes []Envelope) (bool, int64) {
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

func consumePages(pages int, start, end time.Time, cfg Config, f func([]Envelope) (bool, int64)) ([]time.Duration, bool) {
	var queryDurations []time.Duration

	for i := 0; i < pages; i++ {
		queryDuration, last, count, err := consumeTimeRange(start, end, cfg, f)
		if err != nil {
			log.Println(err)
			continue
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

func consumeTimeRange(start, end time.Time, cfg Config, f func([]Envelope) (bool, int64)) (time.Duration, time.Time, int, error) {
	url := fmt.Sprintf("%s/%s?starttime=%d&endtime=%d",
		cfg.LogCacheAddr,
		cfg.VCapApp.ApplicationID,
		start.UnixNano(),
		end.UnixNano(),
	)

	queryStart := time.Now()
	resp, err := http.Get(url)
	if err != nil {
		return 0, time.Time{}, 0, err
	}
	queryDuration := time.Since(queryStart)

	if resp.StatusCode != http.StatusOK {
		return 0, time.Time{}, 0, fmt.Errorf("unxpected status code from log-cache: %d", resp.StatusCode)
	}

	var data LogCacheData
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return 0, time.Time{}, 0, fmt.Errorf("failed to decode response body: %s", err)
	}

	found, last := f(data.Envelopes)
	if !found {
		return queryDuration, time.Unix(0, last), len(data.Envelopes), nil
	}

	return queryDuration, time.Time{}, len(data.Envelopes), nil
}

func lookForMsg(msg string, envelopes []Envelope) (bool, int64) {
	var last int64
	for _, e := range envelopes {
		var err error
		last, err = strconv.ParseInt(e.Timestamp, 10, 64)
		if err != nil {
			log.Printf("failed to parse timestamp: %s", err)
			continue
		}

		if e.Log.Payload != "" {
			payload, err := base64.StdEncoding.DecodeString(e.Log.Payload)
			if err != nil {
				log.Printf("failed to base64 decode log payload: %s", err)
				continue
			}

			if string(payload) == msg {
				return true, last
			}
		}
	}

	return false, last
}

func countMessages(prefix string, envelopes []Envelope) int {
	var count int
	for _, e := range envelopes {
		if e.Log.Payload != "" {
			payload, err := base64.StdEncoding.DecodeString(e.Log.Payload)
			if err != nil {
				log.Printf("failed to base64 decode log payload: %s", err)
				continue
			}

			if strings.Contains(string(payload), prefix) {
				count++
			}
		}
	}

	return count
}

package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
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
	Log Log `json:"log"`
}

type Log struct {
	Payload string `json:"payload"`
}

type TestResult struct {
	Latency          time.Duration   `json:"latency"`
	QueryTimes       []time.Duration `json:"queryTime"`
	AverageQueryTime time.Duration   `json:"averageQueryTime"`
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

		var queryTimes []time.Duration
		for j := 0; j < 10; j++ {
			expectedLog := fmt.Sprintf("Test log - %d", rand.Int63())
			logStartTime := time.Now()
			fmt.Println(expectedLog)

			url := fmt.Sprintf("%s/%s",
				cfg.LogCacheAddr,
				cfg.VCapApp.ApplicationID,
			)

			for i := 0; i < 100; i++ {
				queryStart := time.Now()
				resp, err := http.Get(url)
				if err != nil {
					log.Printf("error from log-cache: %s", err)
					continue
				}

				if resp.StatusCode != http.StatusOK {
					log.Printf("unxpected status code from log-cache: %d", resp.StatusCode)
					continue
				}
				queryTimes = append(queryTimes, time.Since(queryStart))

				var data LogCacheData
				err = json.NewDecoder(resp.Body).Decode(&data)
				if err != nil {
					log.Printf("failed to decode response body: %s", err)
					continue
				}

				if !lookForMsg(expectedLog, data.Envelopes) {
					log.Printf("unable to find msg: %s", expectedLog)
					continue
				}

				var totalQueryTimes time.Duration
				for _, qt := range queryTimes {
					totalQueryTimes += qt
				}

				avgQT := int64(totalQueryTimes) / int64(len(queryTimes))

				// Success
				testResults := TestResult{
					Latency:          time.Since(logStartTime),
					QueryTimes:       queryTimes,
					AverageQueryTime: time.Duration(avgQT),
				}

				resultData, err := json.Marshal(testResults)
				if err != nil {
					log.Printf("failed to marshal test results: %s", err)
					return
				}

				w.Write(resultData)
				return
			}
		}
	})
}

func lookForMsg(msg string, envelopes []Envelope) bool {
	for _, e := range envelopes {
		if e.Log.Payload != "" {
			payload, err := base64.StdEncoding.DecodeString(e.Log.Payload)
			if err != nil {
				log.Printf("failed to base64 decode log payload: %s", err)
				continue
			}

			if string(payload) == msg {
				return true
			}
		}
	}

	return false
}

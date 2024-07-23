package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"code.cloudfoundry.org/go-log-cache/v3/rpc/logcache_v1"
	"code.cloudfoundry.org/go-loggregator/v10/rpc/loggregator_v2"

	client "code.cloudfoundry.org/go-log-cache/v3"
	uuid "github.com/nu7hatch/gouuid"
)

const (
	defaultSampleSize = 10
	messagePrefix     = "loggregator-latency-test-"
)

func main() {
	log.SetOutput(os.Stdout)
	log.Print("starting")

	addr, token, location, err := input()
	if err != nil {
		log.Fatal(err)
	}

	mux := &http.ServeMux{}
	mux.Handle("/", &healthHandler{})
	mux.Handle("/latency", &latencyHandler{
		location: location,
		token:    token,
		mutex:    &sync.Mutex{},
	})

	server := &http.Server{
		Addr:           addr,
		Handler:        mux,
		MaxHeaderBytes: 1 << 20,
	}
	log.Print("listening on " + addr)
	log.Fatal(server.ListenAndServe())
}

func input() (addr, token string, location *url.URL, err error) {
	targetURL := os.Getenv("TARGET_URL")
	if targetURL == "" {
		return "", "", nil, errors.New("empty target url")
	}

	token = os.Getenv("TOKEN")
	if token == "" {
		return "", "", nil, errors.New("empty token")
	}

	port := os.Getenv("PORT")
	if port == "" {
		return "", "", nil, errors.New("empty port")
	}
	addr = ":" + port

	location, err = url.Parse(targetURL)
	if err != nil {
		return "", "", nil, fmt.Errorf("invalid target url: %s", err)
	}

	return addr, token, location, nil
}

func appID() (string, error) {
	appJSON := []byte(os.Getenv("VCAP_APPLICATION"))
	var appData map[string]interface{}
	err := json.Unmarshal(appJSON, &appData)
	if err != nil {
		return "", err
	}
	appID, ok := appData["application_id"].(string)
	if !ok {
		return "", errors.New("can not type assert app id")
	}
	return appID, nil
}

type healthHandler struct{}

func (h *healthHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

type latencyHandler struct {
	location *url.URL
	token    string

	sendTimes map[string]time.Time

	mutex *sync.Mutex
}

func (h *latencyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	sampleSize := sampleSize(r)

	results := h.executeLatencyTest(sampleSize)
	resultBytes, err := json.Marshal(results)
	if err != nil {
		panic(err)
	}

	_, _ = w.Write(resultBytes)
}

func sampleSize(r *http.Request) int {
	samplesQuery := r.URL.Query().Get("samples")
	if samplesQuery == "" {
		return defaultSampleSize
	}
	sampleSize, err := strconv.Atoi(samplesQuery)
	if err != nil || sampleSize < 1 {
		return defaultSampleSize
	}
	return sampleSize
}

type testResults struct {
	AvgSeconds   float64 `json:"avg_seconds"`
	MaxSeconds   float64 `json:"max_seconds"`
	LogsReceived int     `json:"logs_received"`
	LogsExpected int     `json:"logs_expected"`
}

func (h *latencyHandler) executeLatencyTest(sampleQuantity int) testResults {
	h.mutex.Lock()
	h.sendTimes = make(map[string]time.Time)
	h.mutex.Unlock()

	startTime := time.Now()

	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()

	go func() {
		defer func() {
			println("****************** cancelling due to timeout ******************")
			cancelFunc()
		}()

		h.printLogs(sampleQuantity)

		testTimeout := time.Now().Add(30 * time.Second)
		ticker := time.NewTicker(time.Second)
		for range ticker.C {
			fmt.Println("heartbeat")
			if time.Now().After(testTimeout) {
				return
			}
		}
	}()

	results := h.walk(ctx, sampleQuantity, startTime)
	return computeTestResults(results, sampleQuantity)
}

func (h *latencyHandler) walk(ctx context.Context, sampleQuantity int, startTime time.Time) map[string]time.Duration {
	results := make(map[string]time.Duration)

	appID, _ := appID()
	lcClient := client.NewClient(h.location.String(),
		client.WithHTTPClient(&httpClient{
			AuthToken: h.token,
		}))

	visitor := h.visitor(results, sampleQuantity)
	for {
		fmt.Printf("\n%d results recorded so far\n", len(results))

		select {
		case <-ctx.Done():
			return results
		default:
			client.Walk(ctx, appID, visitor, lcClient.Read,
				client.WithWalkEnvelopeTypes(logcache_v1.EnvelopeType_LOG),
				client.WithWalkLimit(1000),
				client.WithWalkLogger(log.New(os.Stdout, "walk: ", 0)),
				client.WithWalkBackoff(client.NewAlwaysRetryBackoff(100*time.Millisecond)),
				client.WithWalkStartTime(startTime.Add(-time.Second)),
				client.WithWalkDelay(time.Nanosecond),
			)

			if len(results) == sampleQuantity {
				return results
			}
		}
	}
}

func (h *latencyHandler) visitor(results map[string]time.Duration, sampleQuantity int) client.Visitor {
	return func(envelopes []*loggregator_v2.Envelope) bool {
		for _, envelope := range envelopes {
			end := time.Now()

			switch envelope.GetMessage().(type) {
			case *loggregator_v2.Envelope_Log:
				message := string(envelope.GetLog().GetPayload())
				if strings.Contains(message, messagePrefix) {
					h.recordResult(results, message, end)
				}
			default:
				continue
			}

			if len(results) == sampleQuantity {
				return false
			}
		}

		return !h.pastDoneSendingTime(results, envelopes[len(envelopes)-1].GetTimestamp())
	}
}

func (h *latencyHandler) pastDoneSendingTime(results map[string]time.Duration, lastEnvelopeTimestamp int64) bool {
	h.mutex.Lock()
	doneSendingTime, ok := h.sendTimes["done_sending"]
	h.mutex.Unlock()

	return ok && time.Unix(0, lastEnvelopeTimestamp).After(doneSendingTime.Add(5*time.Second))
}

func (h *latencyHandler) recordResult(results map[string]time.Duration, message string, end time.Time) {
	_, alreadyReceived := results[message]
	if alreadyReceived {
		return
	}

	h.mutex.Lock()
	start, ok := h.sendTimes[message]
	h.mutex.Unlock()

	if ok {
		results[message] = end.Sub(start)
	}
}

func (h *latencyHandler) printLogs(sampleQuantity int) {
	defer println("done printing log messages")
	for i := 0; i < sampleQuantity; i++ {
		sampleMessage := generateRandomMessage()
		h.mutex.Lock()
		h.sendTimes[sampleMessage] = time.Now()
		fmt.Println(sampleMessage)
		h.mutex.Unlock()
		time.Sleep(time.Millisecond)
	}

	h.mutex.Lock()
	h.sendTimes["done_sending"] = time.Now()
	h.mutex.Unlock()
}

func computeTestResults(results map[string]time.Duration, sampleQuantity int) testResults {
	r := testResults{
		AvgSeconds:   -1,
		MaxSeconds:   -1,
		LogsReceived: len(results),
		LogsExpected: sampleQuantity,
	}
	if len(results) == 0 {
		return r
	}

	var totalDuration, maxDuration time.Duration
	for _, d := range results {
		totalDuration += d
		if d > maxDuration {
			maxDuration = d
		}
	}

	avg := totalDuration / time.Duration(len(results))
	r.AvgSeconds = avg.Seconds()
	r.MaxSeconds = maxDuration.Seconds()
	return r
}

func generateRandomMessage() string {
	id, _ := uuid.NewV4()
	return fmt.Sprint(messagePrefix, id.String())
}

type httpClient struct {
	AuthToken string
}

func (c *httpClient) Do(req *http.Request) (*http.Response, error) {
	req.Header["Authorization"] = []string{c.AuthToken}

	response, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Printf("read failed: %s", err)
		return response, err
	}

	if response.StatusCode != http.StatusOK {
		fmt.Printf("read failed: %d\n", response.StatusCode)
	}
	return response, err
}

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"code.cloudfoundry.org/rfc5424"
)

func main() {
	envInterval := os.Getenv("INTERVAL")
	interval, err := time.ParseDuration(envInterval)
	if err != nil {
		if envInterval != "" {
			log.Fatal(err)
		}

		interval = time.Second
	}

	handler := NewSyslog()
	if os.Getenv("COUNTER_URL") != "" {
		go handler.reportCount(interval)
	}

	fmt.Println(http.ListenAndServe(fmt.Sprintf(":%s", os.Getenv("PORT")), handler))
}

type Counter struct {
	primeCount uint64
	msgCount   uint64
}

// populate handler with json from log message
type Handler struct {
	mu       sync.Mutex
	counters map[string]*Counter
}

func NewSyslog() *Handler {
	return &Handler{
		counters: make(map[string]*Counter),
	}
}

func (h *Handler) reportCount(interval time.Duration) {
	url := os.Getenv("COUNTER_URL") + "/set/"
	if url == "" {
		log.Fatalf("Missing COUNTER_URL environment variable")
	}

	for range time.Tick(interval) {
		payload, err := json.Marshal(h.fetchCounters())
		if err != nil {
			log.Panicf("Failed to marshal counters: %s", err)
		}

		log.Printf("Posting %s", string(payload))
		resp, err := http.Post(url, "application/json", bytes.NewReader(payload))

		if err != nil {
			log.Printf("Failed to write count: %s", err)
			continue
		}

		if resp.StatusCode != http.StatusOK {
			log.Printf("Failed to write count: expected 200 got %d", resp.StatusCode)
		}
	}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close() //nolint:errcheck

	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("Failed to read body: %s", err)
		return
	}

	if len(body) < 1 {
		w.WriteHeader(http.StatusBadRequest)
		log.Print("Empty body")
		return
	}

	// NOTE: This is not a log line. We want the messages printed to STDOUT.
	fmt.Printf("Received: %s", string(body))

	msg := rfc5424.Message{}
	err = msg.UnmarshalBinary(body)
	if err != nil {
		log.Printf("Failed to unmarshal (via RFC-5424) message: %s", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// NOTE: Only one of "primeCount" or "msgCount" will be set.
	var msgCounts messageCount
	err = json.Unmarshal(msg.Message, &msgCounts)
	if err != nil {
		log.Printf("Failed to unmarshal (via JSON) message (%s): %s", string(msg.Message), err)
		return
	}

	h.mu.Lock()
	defer h.mu.Unlock()

	c, ok := h.counters[msgCounts.ID]
	if !ok {
		c = new(Counter)
		h.counters[msgCounts.ID] = c
	}

	c.primeCount += msgCounts.PrimeCount
	c.msgCount += msgCounts.MsgCount
}

type messageCount struct {
	ID         string `json:"id"`
	PrimeCount uint64 `json:"primeCount"`
	MsgCount   uint64 `json:"msgCount"`
}

// fetchCounters extracts the current counter list and returns a slice of
// results in a thread safe manner.
func (h *Handler) fetchCounters() []messageCount {
	h.mu.Lock()
	defer h.mu.Unlock()
	var results []messageCount
	for k, v := range h.counters {
		results = append(results, messageCount{
			ID:         k,
			PrimeCount: v.primeCount,
			MsgCount:   v.msgCount,
		})
	}
	return results
}

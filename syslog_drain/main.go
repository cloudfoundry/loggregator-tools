package main

// TODO: consider backfilling tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"sync"
	"time"

	"code.cloudfoundry.org/rfc5424"
)

var (
	mu sync.Mutex
	// TODO: Has a memory leak. Never deletes entries and constantly adding
	// new and unique keys.
	counters map[string]*Counter
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

	l, err := net.Listen("tcp4", fmt.Sprintf(":%s", os.Getenv("PORT")))
	if err != nil {
		log.Fatal(err)
	}
	defer l.Close()
	log.Print("Listening on " + os.Getenv("PORT"))

	counters = make(map[string]*Counter)
	if os.Getenv("COUNTER_URL") != "" {
		go reportCounts(interval)
	}

	for {
		conn, err := l.Accept()
		log.Printf("Accepted connection")
		if err != nil {
			log.Printf("Error accepting: %s", err)
			continue
		}

		go handle(conn)
	}
}

func reportCounts(interval time.Duration) {
	url := os.Getenv("COUNTER_URL") + "/set/"
	if url == "" {
		log.Fatalf("Missing COUNTER_URL environment variable")
	}

	for range time.Tick(interval) {
		payload, err := json.Marshal(fetchCounters())
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

func handle(conn net.Conn) {
	defer conn.Close()

	var msg rfc5424.Message
	for {
		_, err := msg.ReadFrom(conn)
		if err != nil {
			if err == io.EOF {
				return
			}
			log.Printf("ReadFrom err: %s", err)
			return
		}

		if bytes.Contains(msg.Message, []byte("HTTP")) {
			// TODO: We probably won't ever get this far because a HTTP
			// request won't parse as a RFC-5424.
			continue
		}

		// NOTE: This is not a log line. We want the message to be written
		// to stdout.
		fmt.Printf("Msg %s", string(msg.Message))

		var msgCounts messageCount
		err = json.Unmarshal(msg.Message, &msgCounts)
		if err != nil {
			log.Printf("Failed to unmarshal (via JSON) message (%s): %s", string(msg.Message), err)
			continue
		}

		mu.Lock()
		c, ok := counters[msgCounts.ID]
		if !ok {
			c = new(Counter)
			counters[msgCounts.ID] = c
		}
		c.primeCount += msgCounts.PrimeCount
		c.msgCount += msgCounts.MsgCount
		mu.Unlock()
	}
}

type Counter struct {
	primeCount uint64
	msgCount   uint64
}

type messageCount struct {
	ID         string `json:"id"`
	PrimeCount uint64 `json:"primeCount"`
	MsgCount   uint64 `json:"msgCount"`
}

// fetchCounters extracts the current counter list and returns a slice of
// results in a thread safe manner.
func fetchCounters() []messageCount {
	mu.Lock()
	defer mu.Unlock()
	var results []messageCount
	for k, v := range counters {
		results = append(results, messageCount{
			ID:         k,
			PrimeCount: v.primeCount,
			MsgCount:   v.msgCount,
		})
	}
	return results
}

package web

import (
	"container/ring"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
)

type Counter struct {
	mu     sync.RWMutex
	counts *ring.Ring
}

func NewCounter(n int) Counter {
	return Counter{
		counts: ring.New(n),
	}
}

func (c *Counter) SetHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close() //nolint:errcheck

	var counts []messageCount
	if err := json.NewDecoder(r.Body).Decode(&counts); err != nil {
		log.Printf("Failed to unmarshal JSON request body: %s", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	for _, m := range counts {
		c.counts.Value = m
		c.counts = c.counts.Next()
	}
}

func (c *Counter) GetHandler(w http.ResponseWriter, r *http.Request) {
	id := getID(r)
	var count uint64
	c.mu.RLock()
	c.counts.Do(func(p interface{}) {
		if mc, ok := p.(messageCount); ok {
			if mc.ID == id {
				count = mc.MsgCount
			}
		}
	})
	c.mu.RUnlock()

	_, err := fmt.Fprint(w, count)
	if err != nil {
		fmt.Println("could not write response:", err)
	}
}

func (c *Counter) GetPrimeHandler(w http.ResponseWriter, r *http.Request) {
	id := getID(r)
	var count uint64
	c.mu.RLock()
	c.counts.Do(func(p interface{}) {
		if mc, ok := p.(messageCount); ok {
			if mc.ID == id {
				count = mc.PrimeCount
			}
		}
	})
	c.mu.RUnlock()

	_, err := fmt.Fprint(w, count)
	if err != nil {
		fmt.Println("could not write response:", err)
	}
}

type messageCount struct {
	ID         string `json:"id"`
	PrimeCount uint64 `json:"primeCount"`
	MsgCount   uint64 `json:"msgCount"`
}

func getID(r *http.Request) string {
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 1 {
		return ""
	}

	return parts[len(parts)-1]
}

package main

// TODO: consider backfilling tests

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
)

type counter struct {
	prime uint64
	msg   uint64
}

var (
	mu sync.Mutex

	// TODO This map leaks memory.
	counters map[string]*counter
)

func main() {
	verbose := os.Getenv("VERBOSE")
	if verbose != "true" {
		log.SetOutput(ioutil.Discard)
	}

	counters = make(map[string]*counter)

	http.Handle("/get/", http.HandlerFunc(getCountHandler))
	http.Handle("/get-prime/", http.HandlerFunc(getPrimeCountHandler))
	// TODO: change set to take in a json payload vs a single ID/counters pair
	http.Handle("/set/", http.HandlerFunc(setCountHandler))
	http.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("404 - %+v", r.URL)
		w.WriteHeader(http.StatusNotFound)
	}))

	addr := fmt.Sprintf(":%s", os.Getenv("PORT"))
	log.Print("Listening on " + os.Getenv("PORT"))
	log.Println(http.ListenAndServe(addr, nil))
}

func getID(r *http.Request) string {
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 1 {
		return ""
	}

	log.Printf("Using ID %s", parts[len(parts)-1])
	return parts[len(parts)-1]
}

func getCounter(id string) counter {
	mu.Lock()
	defer mu.Unlock()
	c := counters[id]
	if c == nil {
		c = &counter{}
		counters[id] = c
	}
	return *c
}

func setCounter(id string, prime, msg uint64) {
	mu.Lock()
	defer mu.Unlock()
	c := counters[id]
	if c == nil {
		c = &counter{}
		counters[id] = c
	}
	c.prime = prime
	c.msg = msg
}

func getCountHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("GET /get")

	counter := getCounter(getID(r))

	w.Write([]byte(fmt.Sprint(counter.msg)))
}

func getPrimeCountHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("GET /get-prime")

	counter := getCounter(getID(r))

	w.Write([]byte(fmt.Sprint(counter.prime)))
}

type messageCount struct {
	ID         string `json:"id"`
	PrimeCount uint64 `json:"primeCount"`
	MsgCount   uint64 `json:"msgCount"`
}

// setCountHandler expects a JSON payload that adheres to the following
// structure:
// [{
//   id: string,
//   primeCount: number,
//   msgCount: number,
// }]
func setCountHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	log.Println("POST /set")

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Failed to read request body: %s", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	log.Printf("Body: %s", string(body))

	var counts []messageCount
	err = json.Unmarshal(body, &counts)
	if err != nil {
		log.Printf("Failed to unmarshal JSON request body: %s", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	for _, m := range counts {
		setCounter(m.ID, m.PrimeCount, m.MsgCount)
	}
}

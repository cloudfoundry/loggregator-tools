package main

// TODO: consider backfilling tests

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"code.cloudfoundry.org/loggregator-tools/counter/web"
)

func main() {
	verbose := os.Getenv("VERBOSE")
	if verbose != "true" {
		log.SetOutput(ioutil.Discard)
	}

	counter := web.NewCounter(100)

	http.HandleFunc("/get/", counter.GetHandler)
	http.HandleFunc("/get-prime/", counter.GetPrimeHandler)
	http.HandleFunc("/set/", counter.SetHandler)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("404 - %+v", r.URL)
		w.WriteHeader(http.StatusNotFound)
	})

	addr := fmt.Sprintf(":%s", os.Getenv("PORT"))
	log.Print("Listening on " + os.Getenv("PORT"))
	log.Println(http.ListenAndServe(addr, nil))
}

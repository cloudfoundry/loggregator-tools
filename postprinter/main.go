package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func main() {
	log.Print("Starting Post Printer...")
	defer log.Print("Closing Post Printer.")

	port := os.Getenv("PORT")

	var logRequests bool
	logRequests = os.Getenv("SKIP_REQUEST_LOGGING") == "" || os.Getenv("SKIP_REQUEST_LOGGING") != "true"

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		if logRequests {
			log.Printf("Request: %+v", r)
		}
		data, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Printf("Error reading body: %s", err)
			return
		}

		log.Printf("Body: %s", data)

		if r.Method == "POST" {
			w.WriteHeader(http.StatusCreated)
			return
		}

		w.WriteHeader(http.StatusOK)
	})

	http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
}

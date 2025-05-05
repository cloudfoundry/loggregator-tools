package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

func main() {
	log.Print("Starting Post Printer...")
	defer log.Print("Closing Post Printer.")

	port := os.Getenv("PORT")

	logRequests := os.Getenv("SKIP_REQUEST_LOGGING") == "" || os.Getenv("SKIP_REQUEST_LOGGING") != "true"

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close() //nolint:errcheck

		if logRequests {
			log.Printf("Request: %+v", r)
		}

		data, err := io.ReadAll(r.Body)
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

	fmt.Println(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}

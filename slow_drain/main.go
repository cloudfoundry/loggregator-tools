package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func main() {
	log.Print("Starting slow drain...")
	defer log.Print("Closing slow drain.")

	port := os.Getenv("PORT")

	var logRequests bool
	logRequests = os.Getenv("SKIP_REQUEST_LOGGING") == "" || os.Getenv("SKIP_REQUEST_LOGGING") != "true"

	blockChan := make(chan int)
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
		<-blockChan
	})

	http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
}

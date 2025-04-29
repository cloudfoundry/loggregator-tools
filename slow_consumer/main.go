package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/cloudfoundry/noaa/v2/consumer"
)

var (
	loggrAddr    = os.Getenv("LOGGREGATOR_ADDR")
	authToken    = os.Getenv("CF_ACCESS_TOKEN")
	readInterval = os.Getenv("READ_INTERVAL")
	shardID      = os.Getenv("SHARD_ID")
)

func main() {
	var missing []string
	if loggrAddr == "" {
		missing = append(missing, "LOGGREGATOR_ADDR")
	}

	if authToken == "" {
		missing = append(missing, "CF_ACCESS_TOKEN")
	}

	if len(missing) > 0 {
		log.Fatalf("Missing required environment variables: %s", strings.Join(missing, ", "))
	}

	if shardID == "" {
		shardID = fmt.Sprintf("shard-id-%d", rand.Int63())
	}

	if readInterval == "" {
		readInterval = "10s"
	}

	ri, err := time.ParseDuration(readInterval)
	if err != nil {
		log.Fatalf("failed to parse read interval: %s", err)
	}

	log.Println("Starting slow consumer")
	cnsmr := consumer.New(loggrAddr, &tls.Config{InsecureSkipVerify: true}, nil)
	msgChan, errorChan := cnsmr.Firehose(shardID, authToken)

	go func() {
		for err := range errorChan {
			log.Printf("Error while reading from firehose: %s", err)
		}
	}()

	for range msgChan {
		log.Printf("Received message, waiting for %s to read the next", ri)
		time.Sleep(ri)
	}
}

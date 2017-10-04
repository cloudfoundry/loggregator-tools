package main

import (
	"log"
	"os"

	"code.cloudfoundry.org/loggregator-tools/syslog_to_datadog/app"
)

func main() {
	datadogAPIKey := os.Getenv("DATADOG_API_KEY")
	if datadogAPIKey == "" {
		log.Fatal("missing required environment variable DATADOG_API_KEY")
	}

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("missing required environment variable PORT")
	}

	s := app.NewServer(":"+port, datadogAPIKey)
	s.Run()
}

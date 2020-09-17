// tool for generating log load

// usage: gets settings from application name
// If pushed as just `constlogger`, emits 1000 logs per second
// If pushed as, e.g., `constlogger-100`, emits 100 logs per second
// If pushed as, e.g., `constlogger-100-50`, emits 100 logs per second that are each 50 bytes long (including the newline)

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"go.uber.org/ratelimit"
)

type VcapApplication struct {
	ApplicationName string `json:"application_name"`
}

func main() {
	var app VcapApplication
	err := json.Unmarshal([]byte(os.Getenv("VCAP_APPLICATION")), &app)
	if err != nil {
		log.Panicf("failed to parse VCAP_APPLICATION: %q", err)
	}

	rate := 1000
	minLength := 0

	parts := strings.Split(app.ApplicationName, "-")

	if len(parts) > 1 {
		rate, err = strconv.Atoi(parts[1])

		if err != nil {
			log.Panicf("failed to parse rate from application name (constlogger-RATE): %q", err)
		}
	}

	if len(parts) > 2 {
		minLength, err = strconv.Atoi(parts[2])

		if err != nil {
			log.Panicf("failed to parse minimum log length from application name (constlogger-RATE-LENGTH): %q", err)
		}
	}

	go func() {
		err := http.ListenAndServe(":"+os.Getenv("PORT"), nil)
		if err != nil {
			panic(err)
		}
	}()

	fmt.Printf("logging %d msgs/sec that are at least %d bytes long\n", rate, minLength)

	total := 0
	rl := ratelimit.New(rate)
	for {
		total += 1
		logMsg := fmt.Sprintf("msg %d", total)
		padLen := minLength - len(logMsg) - 2
		if padLen < 0 {
			padLen = 0
		}
		fmt.Printf("%s %s\n", logMsg, strings.Repeat("-", padLen))
		rl.Take()
	}
}

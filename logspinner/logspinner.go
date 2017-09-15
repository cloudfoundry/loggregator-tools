// tool for generating log load

// usage: curl <endpoint>?cycles=100&delay=1ms&text=time2
// delay is duration format (https://golang.org/pkg/time/#ParseDuration)
// defaults: 10 cycles, 1 second, "LogSpinner Log Message"

package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

func main() {
	http.HandleFunc("/", rootResponse)
	fmt.Println("listening...")
	err := http.ListenAndServe(":"+os.Getenv("PORT"), nil)
	if err != nil {
		panic(err)
	}
}

func rootResponse(res http.ResponseWriter, req *http.Request) {
	cycleCount, err := strconv.Atoi(req.FormValue("cycles"))
	if cycleCount == 0 || err != nil {
		cycleCount = 10
	}

	delay, err := time.ParseDuration(req.FormValue("delay"))
	if err != nil {
		delay = 1000 * time.Millisecond
	}

	id := req.FormValue("id")
	if id == "" {
		id = fmt.Sprintf("%d", time.Now().UnixNano())
	}

	mode := "msgCount"
	isPrimer := (req.FormValue("primer"))
	if isPrimer == "true" {
		mode = "primeCount"
	}

	go outputLog(cycleCount, delay, id, mode)

	fmt.Fprintf(res, "cycles %d, delay %s, id %s, mode %s\n", cycleCount, delay, id, mode)
}

func outputLog(cycleCount int, delay time.Duration, id, mode string) {
	now := time.Now()
	for i := 0; i < cycleCount; i++ {
		payload := fmt.Sprintf(`{"id":%q,"cycles":%d,"delay":%q,%q:1,"iteration":%d}`, id, cycleCount, delay, mode, i+1)
		fmt.Println(payload)
		time.Sleep(delay)
	}
	done := time.Now()
	diff := done.Sub(now)

	rate := float64(cycleCount) / diff.Seconds()
	log.Printf("Duration %s TotalSent %d Rate %f", diff.String(), cycleCount, rate)
}

// tool for generating log load

// usage: curl <endpoint>?rate=100&duration=1ms&text=time2
// duration is duration format (https://golang.org/pkg/time/#ParseDuration)
// defaults: 100 logs per second, 1 second, "LogSpinner Log Message"

package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"go.uber.org/ratelimit"
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

	rate, err := strconv.Atoi(req.FormValue("rate"))
	if rate == 0 || err != nil {
		rate = 100
	}

	duration, err := time.ParseDuration(req.FormValue("duration"))
	if err != nil {
		duration = time.Second
	}

	logText := (req.FormValue("text"))
	if logText == "" {
		logText = "LogSpinner Log Message"
	}

	go outputLog(rate, duration, logText)

	_, err = fmt.Fprintf(res, "rate %d, duration %s, text %s\n", rate, duration, logText)
	if err != nil {
		fmt.Println("error writing response:", err)
	}
}

func outputLog(rate int, duration time.Duration, logText string) {
	end := time.Now().Add(duration)
	rl := ratelimit.New(rate)
	now := time.Now()
	total := 0
	for time.Now().Before(end) {
		total += 1
		fmt.Printf("msg %d %s\n", total, logText)
		rl.Take()
	}
	done := time.Now()
	diff := done.Sub(now)

	actualRate := float64(total) / diff.Seconds()
	fmt.Printf("Duration %s TotalSent %d Rate %f \n", diff.String(), total, actualRate)

}

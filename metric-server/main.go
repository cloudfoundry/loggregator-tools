package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
)

const metrics = `
# HELP %s_node_timex_pps_calibration_total Pulse per second count of calibration intervals.
# TYPE %s_node_timex_pps_calibration_total counter
%s_node_timex_pps_calibration_total 1
`

func main() {
	http.HandleFunc("/", helloWorldHandler)
	http.HandleFunc("/metrics", metricHandler)

	go http.ListenAndServe(":"+os.Getenv("PORT"), nil) //nolint:errcheck
	fmt.Println(http.ListenAndServe(":8081", nil))
}

func helloWorldHandler(w http.ResponseWriter, r *http.Request) {
	_, err := fmt.Fprint(w, "Hello World")
	if err != nil {
		fmt.Println("error writing response:", err)
	}
}

func metricHandler(w http.ResponseWriter, r *http.Request) {
	var instances = []string{"a", "b", "c", "d"}
	index := os.Getenv("CF_INSTANCE_INDEX")

	i, err := strconv.Atoi(index)
	if err != nil {
		panic(err)
	}

	_, err = fmt.Fprintf(w, metrics, instances[i%4], instances[i%4], instances[i%4])
	if err != nil {
		fmt.Println("error writing response:", err)
	}
}

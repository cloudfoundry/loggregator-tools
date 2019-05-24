package main

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/google/uuid"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rcrowley/go-metrics"
)

const RFC5424TimeOffsetNum = "2006-01-02T15:04:05.999999-07:00"

var egress = promauto.NewCounter(prometheus.CounterOpts{
	Name: "biglogger_egress",
	Help: "The total number of egressed logs",
})

var dropped = promauto.NewCounter(prometheus.CounterOpts{
	Name: "biglogger_dropped",
	Help: "The total number of dropped logs",
})

var reset = promauto.NewCounter(prometheus.CounterOpts{
	Name: "biglogger_connections_reset",
	Help: "The total number of reset connections",
})

var egressRateMinute = promauto.NewGauge(prometheus.GaugeOpts{
	Name: "biglogger_egress_rate_minute",
	Help: "egress rate",
})

var egressDuration = promauto.NewHistogram(prometheus.HistogramOpts{
	Name: "biglogger_egress_duration",
	Help: "egress rate",
})

var egressRate = metrics.NewMeter()

func main() {
	logsPerSecondPerEmitter := os.Getenv("LOGS_PER_SECOND_PER_EMITTER")
	ipsString := os.Getenv("IPS")
	numEmitters := os.Getenv("NUM_EMITTERS")

	logsPerSecond, err := strconv.Atoi(logsPerSecondPerEmitter)
	if err != nil {
		log.Panic("failed to convert logs per second")
	}

	ips := strings.Split(ipsString, ",")
	ne, err := strconv.Atoi(numEmitters)
	if err != nil {
		log.Panic("failed to convert num emitters")
	}

	log.Print("Starting writers")
	for i := 0; i < ne; i++ {
		ip := ips[i%len(ips)]
		go writeLogs(logsPerSecond, ip)
		log.Printf("Started writer for ip: %s", ip)
	}

	http.Handle("/metrics", promhttp.Handler())
	go func() {
		err = http.ListenAndServe(fmt.Sprintf(":%s", os.Getenv("PORT")), nil)
		if err != nil {
			log.Fatalf("Failed to start metrics server: %s", err)
		}
	}()

	go func() {
		t := time.NewTicker(time.Second)
		for range t.C {
			egressRateMinute.Set(egressRate.Snapshot().Rate1())
		}
	}()

	sigs := make(chan os.Signal)
	signal.Notify(sigs, syscall.SIGTERM, syscall.SIGINT)
	<-sigs

	log.Println("EXITING")
}

func writeLogs(logsPerSecond int, ip string) {
	t := time.NewTicker(time.Second)
	guid := uuid.New()
	conn := connect(ip)
	defer conn.Close()

	for range t.C {
		conn = emitBatch(logsPerSecond, guid, ip, conn)
	}
}

func emitBatch(batchSize int, guid uuid.UUID, ip string, conn net.Conn) net.Conn {
	for i := 0; i < batchSize; i++ {
		msg := fmt.Sprintf("<14>1 %s test-hostname %s [MY-TASK/2] - - just a test\n", time.Now().Format(RFC5424TimeOffsetNum), guid)

		var err error
		start := time.Now()
		conn, err = writeWithRetry(conn, ip, fmt.Sprintf("%d %s", len([]byte(msg)), msg))
		end := time.Now()

		if err != nil {
			log.Printf("Error writing to log cache: %s\n", err.Error())
			dropped.Inc()
		} else {
			egressDuration.Observe(end.Sub(start).Seconds())
			egress.Inc()
			egressRate.Mark(1)
		}
	}

	return conn
}

func writeWithRetry(conn net.Conn, ip, msg string) (net.Conn, error) {
	log.Println(msg)

	_, err := conn.Write([]byte(msg))
	if err != nil {
		conn.Close()
		conn = connect(ip)

		if opErr, ok := err.(*net.OpError); ok {
			if syscallErr, ok := opErr.Err.(*os.SyscallError); ok {
				if syscallErr.Err == syscall.ECONNRESET {
					reset.Inc()
					_, err = conn.Write([]byte(msg))
				}
			}
		}
	}

	return conn, err
}

func connect(ip string) net.Conn {
	for {
		conn, err := net.Dial("tcp", ip+":8888")
		if err != nil {
			log.Printf("failed connect to endpoint %s", ip)
			time.Sleep(100 * time.Millisecond)
		} else {
			return conn
		}
	}
}

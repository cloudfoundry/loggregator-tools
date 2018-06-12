package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	envstruct "code.cloudfoundry.org/go-envstruct"
	loggregator "code.cloudfoundry.org/go-loggregator"
	"code.cloudfoundry.org/loggregator-tools/syslog-nozzle/app"
	"code.cloudfoundry.org/loggregator/plumbing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	droppedCounter = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "syslog_nozzle_dropped",
		Help: "The count of messages dropped when writing to syslog.",
	})
	egressedCounter = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "syslog_nozzle_egress",
		Help: "The count of messages written to syslog.",
	})
	ignoredCounter = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "syslog_nozzle_ignored",
		Help: "The count of envelopes ignored because of namespace mismatch",
	})
	convertCounter = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "syslog_nozzle_convert_err",
		Help: "The count of envelopes ignored because it couldn't be converted to syslog",
	})
)

func init() {
	prometheus.MustRegister(droppedCounter)
	prometheus.MustRegister(egressedCounter)
	prometheus.MustRegister(ignoredCounter)
	prometheus.MustRegister(convertCounter)
}

func main() {
	var conf app.Config
	err := envstruct.Load(&conf)
	if err != nil {
		log.Fatal(err)
	}
	if conf.ShardID == "" {
		conf.ShardID = generateShardID()
	}

	err = envstruct.WriteReport(&conf)
	if err != nil {
		log.Fatal(err)
	}

	tlsConfig, err := plumbing.NewClientMutualTLSConfig(
		conf.LogsProviderTLS.Cert,
		conf.LogsProviderTLS.Key,
		conf.LogsProviderTLS.CA,
		"reverselogproxy",
	)
	if err != nil {
		log.Fatal(err)
	}
	sc := loggregator.NewEnvelopeStreamConnector(
		conf.LogsProviderAddr,
		tlsConfig,
	)

	nozzle := app.NewNozzle(
		sc,
		conf.Destination,
		conf.ShardID,
		app.WithNamespace(conf.Namespace),
		app.WithEgressedCounter(egressedCounter),
		app.WithDroppedCounter(droppedCounter),
		app.WithIgnoredEnvelopeCounter(ignoredCounter),
		app.WithConversionErrorCounter(convertCounter),
	)
	go nozzle.Start()

	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())
	srv := &http.Server{
		Addr:         conf.MetricsAddr,
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	log.Fatal(srv.ListenAndServe())
}

func generateShardID() string {
	return fmt.Sprintf("syslog_nozzle_%s", time.Now().String())
}

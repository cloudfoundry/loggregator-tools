package main

import (
	"log"

	datadog "github.com/zorkian/go-datadog-api"

	envstruct "code.cloudfoundry.org/go-envstruct"
	logcache "code.cloudfoundry.org/go-log-cache"
	"code.cloudfoundry.org/go-loggregator/rpc/loggregator_v2"
)

func main() {
	log.Print("Starting DataDog Accumulator...")
	defer log.Print("Closing DataDog Accumulator.")

	cfg := loadConfig()

	llc := logcache.NewClient(cfg.LogCacheAddr)
	ddc := datadog.NewClient(cfg.DataDogAPIKey, "")

	logcache.Walk(cfg.SourceID, buildDataDogWriter(ddc), llc.Read)
}

func buildDataDogWriter(ddc *datadog.Client) func([]*loggregator_v2.Envelope) bool {
	return func(es []*loggregator_v2.Envelope) bool {
		metrics := make(map[string]datadog.Metric)
		for _, e := range es {
			if e.Gauge == nil {
				continue
			}

			for name, value := range e.GetGauge().Metrics {
				// We plan to take the address of this and therefore can not
				// use name given to us via range.
				name := name
				mType := "gauge"
				ddc.PostMetric([]datadog.Metric{
					{
						Metric: &name,
						Points: [2]float64{float64(e.Timestamp, value)},
						Type:   &mType,
					},
				})
			}
		}
		return true
	}
}

type Config struct {
	LogCacheAddr  string `env:"LOG_CACHE_ADDR,required"`
	SourceID      string `env:"SOURCE_ID,required"`
	DataDogAPIKey string `env:"DATA_DOG_API_KEY,required"`
}

func loadConfig() *Config {
	var cfg Config
	if err := envstruct.Load(&cfg); err != nil {
		log.Fatalf("failed to load config: %s", err)
	}

	return &cfg
}

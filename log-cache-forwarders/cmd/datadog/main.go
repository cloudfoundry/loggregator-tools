package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"

	envstruct "code.cloudfoundry.org/go-envstruct"
	logcache "code.cloudfoundry.org/go-log-cache"
	"code.cloudfoundry.org/go-loggregator/rpc/loggregator_v2"
	"code.cloudfoundry.org/loggregator-tools/log-cache-forwarders/pkg/groupmanager"
	"code.cloudfoundry.org/loggregator-tools/log-cache-forwarders/pkg/sourceidprovider"
	datadog "github.com/zorkian/go-datadog-api"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	cfg := LoadConfig()
	envstruct.WriteReport(&cfg)

	groupClient := logcache.NewShardGroupReaderClient(
		cfg.LogCacheHTTPAddr,
		logcache.WithHTTPClient(newOauth2HTTPClient(cfg)),
	)

	client := logcache.NewClient(
		cfg.LogCacheHTTPAddr,
		logcache.WithHTTPClient(newOauth2HTTPClient(cfg)),
	)

	if cfg.SourceIDWhitelist != "" && cfg.SourceIDBlacklist != "" {
		log.Fatalf("Can't have a whitelist and a blacklist...")
	}

	var provider groupmanager.GroupProvider = sourceidprovider.All(client)

	if cfg.SourceIDBlacklist != "" {
		provider = sourceidprovider.NewRegex(
			true,
			cfg.SourceIDBlacklist,
			client,
		)
	}

	if cfg.SourceIDWhitelist != "" {
		provider = sourceidprovider.NewRegex(
			false,
			cfg.SourceIDWhitelist,
			client,
		)
	}

	groupmanager.Start(
		cfg.LogCacheGroupName,
		time.Tick(30*time.Second),
		provider,
		groupClient,
	)

	ddc := datadog.NewClient(cfg.DatadogAPIKey, "")
	visitor := buildDatadogWriter(ddc, cfg.MetricHost, strings.Split(cfg.DatadogTags, ","))

	reader := groupClient.BuildReader(rand.Uint64())

	logcache.Walk(
		context.Background(),
		cfg.LogCacheGroupName,
		logcache.Visitor(visitor),
		reader,
		logcache.WithWalkStartTime(time.Now()),
		logcache.WithWalkBackoff(logcache.NewAlwaysRetryBackoff(250*time.Millisecond)),
		logcache.WithWalkLogger(log.New(os.Stderr, "", log.LstdFlags)),
	)
}

func newOauth2HTTPClient(cfg Config) *logcache.Oauth2HTTPClient {
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: cfg.SkipCertVerify,
			},
		},
		Timeout: 5 * time.Second,
	}

	return logcache.NewOauth2HTTPClient(
		cfg.UAAAddr,
		cfg.ClientID,
		cfg.ClientSecret,
		logcache.WithOauth2HTTPClient(client),
	)
}

func buildDatadogWriter(ddc *datadog.Client, host string, tags []string) func([]*loggregator_v2.Envelope) bool {
	return func(es []*loggregator_v2.Envelope) bool {
		var metrics []datadog.Metric
		for _, e := range es {
			switch e.Message.(type) {
			case *loggregator_v2.Envelope_Gauge:
				for name, value := range e.GetGauge().Metrics {
					// We plan to take the address of this and therefore can not
					// use name given to us via range.
					name := name
					if e.GetSourceId() != "" {
						name = fmt.Sprintf("%s.%s", e.GetSourceId(), name)
					}

					mType := "gauge"
					metrics = append(metrics, datadog.Metric{
						Metric: &name,
						Points: toDataPoint(e.Timestamp, value.GetValue()),
						Type:   &mType,
						Host:   &host,
						Tags:   tags,
					})
				}
			case *loggregator_v2.Envelope_Counter:
				name := e.GetCounter().GetName()
				if e.GetSourceId() != "" {
					name = fmt.Sprintf("%s.%s", e.GetSourceId(), name)
				}

				mType := "gauge"
				metrics = append(metrics, datadog.Metric{
					Metric: &name,
					Points: toDataPoint(e.Timestamp, float64(e.GetCounter().GetTotal())),
					Type:   &mType,
					Host:   &host,
					Tags:   tags,
				})
			default:
				continue
			}
		}

		if len(metrics) < 1 {
			return true
		}

		err := ddc.PostMetrics(metrics)
		if err != nil {
			log.Printf("failed to write metrics to DataDog: %s", err)
		}
		log.Printf("posted %d metrics", len(metrics))

		return true
	}
}

func toDataPoint(x int64, y float64) []datadog.DataPoint {
	t := time.Unix(0, x)
	tf := float64(t.Unix())
	return []datadog.DataPoint{
		[2]float64{tf, y},
	}
}

package datadog

import (
	"fmt"
	"log"
	"time"

	"code.cloudfoundry.org/go-loggregator/rpc/loggregator_v2"
	datadog "github.com/zorkian/go-datadog-api"
)

type Client interface {
	PostMetrics(m []datadog.Metric) error
}

func Visitor(c Client, host string, tags []string) func(es []*loggregator_v2.Envelope) bool {
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

		if err := c.PostMetrics(metrics); err != nil {
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
		[2]*float64{&tf, &y},
	}
}

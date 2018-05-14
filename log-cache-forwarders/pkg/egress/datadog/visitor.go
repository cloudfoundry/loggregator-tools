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
	PostEvent(e *datadog.Event) (*datadog.Event, error)
}

func Visitor(c Client, host string, tags []string) func(es []*loggregator_v2.Envelope) bool {
	return func(es []*loggregator_v2.Envelope) bool {
		var metrics []datadog.Metric
		var events []*datadog.Event

		for _, e := range es {
			ddtags := append(make([]string, 0), tags...)

			for key, value := range e.GetTags() {
				ddtags = append(ddtags, key+":"+value)
			}

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
						Tags:   ddtags,
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
					Tags:   ddtags,
				})
			case *loggregator_v2.Envelope_Event:
				event := e.GetEvent()
				title := event.GetTitle()
				text := event.GetBody()

				events = append(events, &datadog.Event{
					Title: &title,
					Text:  &text,
					Host:  &host,
					Tags:  ddtags,
				})
			default:
				continue
			}
		}

		if len(metrics) > 0 {
			if err := c.PostMetrics(metrics); err != nil {
				log.Printf("failed to write metrics to DataDog: %s", err)
			} else {
				log.Printf("posted %d metrics", len(metrics))
			}
		}

		if len(events) > 0 {
			successfulSends := 0
			for _, e := range events {
				if _, err := c.PostEvent(e); err != nil {
					log.Printf("failed to write event to DataDog: %s", err)
				} else {
					successfulSends++
				}
			}

			log.Printf("posted %d events", successfulSends)
		}

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

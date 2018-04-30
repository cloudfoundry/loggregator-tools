package datadog_test

import (
	"code.cloudfoundry.org/go-loggregator/rpc/loggregator_v2"
	"code.cloudfoundry.org/loggregator-tools/log-cache-forwarders/pkg/egress/datadog"
	datadogapi "github.com/zorkian/go-datadog-api"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Visitor", func() {
	Context("when the envelope type is Counter", func() {
		It("writes counter metrics to the datadog client", func() {
			ddc := &stubDatadogClient{}
			visitor := datadog.Visitor(ddc, "hostname", []string{"tag-1", "tag-2"})

			cont := visitor([]*loggregator_v2.Envelope{
				{
					Timestamp: 1000000000,
					Message: &loggregator_v2.Envelope_Counter{
						Counter: &loggregator_v2.Counter{
							Name:  "counter-a",
							Total: 123,
						},
					},
				},
				{
					Timestamp: 1000000000,
					Message: &loggregator_v2.Envelope_Counter{
						Counter: &loggregator_v2.Counter{
							Name:  "counter-b",
							Total: 456,
						},
					},
				},
			})

			Expect(cont).To(BeTrue())
			Expect(ddc.metrics).To(HaveLen(2))

			m := ddc.metrics[0]
			Expect(*m.Type).To(Equal("gauge"))
			Expect(*m.Metric).To(Equal("counter-a"))
			Expect(*m.Host).To(Equal("hostname"))
			Expect(m.Tags).To(ConsistOf("tag-1", "tag-2"))

			Expect(m.Points).To(HaveLen(1))

			p := m.Points[0]

			Expect(*p[0]).To(Equal(float64(1)))
			Expect(*p[1]).To(Equal(float64(123)))
		})

		It("metric name includes source id if present", func() {
			ddc := &stubDatadogClient{}
			visitor := datadog.Visitor(ddc, "hostname", []string{})

			visitor([]*loggregator_v2.Envelope{
				{
					Timestamp: 1000000000,
					SourceId:  "counter-id-1",
					Message: &loggregator_v2.Envelope_Counter{
						Counter: &loggregator_v2.Counter{
							Name:  "counter-a",
							Total: 123,
						},
					},
				},
			})
			m := ddc.metrics[0]
			Expect(*m.Metric).To(Equal("counter-id-1.counter-a"))
		})
	})

	Context("when envelopes is empty", func() {
		It("does not post metrics", func() {
			ddc := &stubDatadogClient{}
			visitor := datadog.Visitor(ddc, "hostname", []string{})

			visitor(nil)

			Expect(ddc.postMetricsCalled).To(BeFalse())
		})
	})

	Context("when the envelope type is Gauge", func() {
		It("writes gauge metrics to the datadog client", func() {
			ddc := &stubDatadogClient{}
			visitor := datadog.Visitor(ddc, "hostname", []string{"tag-1", "tag-2"})

			cont := visitor([]*loggregator_v2.Envelope{
				{
					Timestamp: 1000000000,
					Message: &loggregator_v2.Envelope_Gauge{
						Gauge: &loggregator_v2.Gauge{
							Metrics: map[string]*loggregator_v2.GaugeValue{
								"gauge-a": &loggregator_v2.GaugeValue{
									Unit:  "bytes",
									Value: float64(100),
								},
								"gauge-b": &loggregator_v2.GaugeValue{
									Unit:  "bytes",
									Value: float64(100),
								},
							},
						},
					},
				},
				{
					Timestamp: 1000000000,
					Message: &loggregator_v2.Envelope_Gauge{
						Gauge: &loggregator_v2.Gauge{
							Metrics: map[string]*loggregator_v2.GaugeValue{
								"gauge-c": &loggregator_v2.GaugeValue{
									Unit:  "bytes",
									Value: float64(100),
								},
								"gauge-d": &loggregator_v2.GaugeValue{
									Unit:  "bytes",
									Value: float64(100),
								},
							},
						},
					},
				},
			})

			Expect(cont).To(BeTrue())
			Expect(ddc.metrics).To(HaveLen(4))

			var m *datadogapi.Metric
			for _, metric := range ddc.metrics {
				if *metric.Metric == "gauge-a" {
					m = &metric
					break
				}
			}

			Expect(m).ToNot(BeNil())
			Expect(*m.Type).To(Equal("gauge"))
			Expect(*m.Host).To(Equal("hostname"))
			Expect(m.Tags).To(ConsistOf("tag-1", "tag-2"))

			Expect(m.Points).To(HaveLen(1))

			p := m.Points[0]

			Expect(*p[0]).To(Equal(float64(1)))
			Expect(*p[1]).To(Equal(float64(100)))
		})

		It("metric name includes source id if present", func() {
			ddc := &stubDatadogClient{}
			visitor := datadog.Visitor(ddc, "hostname", []string{})

			visitor([]*loggregator_v2.Envelope{
				{
					Timestamp: 1000000000,
					SourceId:  "gauge-id-1",
					Message: &loggregator_v2.Envelope_Gauge{
						Gauge: &loggregator_v2.Gauge{
							Metrics: map[string]*loggregator_v2.GaugeValue{
								"gauge-a": &loggregator_v2.GaugeValue{
									Unit:  "bytes",
									Value: float64(100),
								},
							},
						},
					},
				},
			})

			m := ddc.metrics[0]
			Expect(*m.Metric).To(Equal("gauge-id-1.gauge-a"))
		})
	})
})

type stubDatadogClient struct {
	postMetricsCalled bool
	metrics           []datadogapi.Metric
}

func (s *stubDatadogClient) PostMetrics(m []datadogapi.Metric) error {
	s.postMetricsCalled = true
	s.metrics = append(s.metrics, m...)
	return nil
}

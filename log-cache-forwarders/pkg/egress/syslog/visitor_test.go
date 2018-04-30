package syslog_test

import (
	"errors"

	"code.cloudfoundry.org/go-loggregator/rpc/loggregator_v2"
	"code.cloudfoundry.org/loggregator-tools/log-cache-forwarders/pkg/egress/syslog"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Visitor", func() {
	It("writes out a collection of envelopes, one at a time", func() {
		spyWriter := newSpyWriter()
		spyMetrics := newSpyMetrics()
		v := syslog.NewVisitor(spyWriter, spyMetrics)
		envs := []*loggregator_v2.Envelope{
			{SourceId: "source-1"},
			{SourceId: "source-2"},
			{SourceId: "source-3"},
		}

		cont := v(envs)

		Expect(cont).To(BeTrue())
		Expect(spyWriter.called).To(Equal(3))
		Expect(spyWriter.calledWith).To(Equal(envs))
	})

	It("increments ingress and syslog metric when envelope is written", func() {
		spyWriter := newSpyWriter()
		spyMetrics := newSpyMetrics()
		v := syslog.NewVisitor(spyWriter, spyMetrics)
		envs := []*loggregator_v2.Envelope{
			{SourceId: "source-1"},
			{SourceId: "source-2"},
			{SourceId: "source-3"},
		}

		v(envs)

		Expect(spyMetrics.counters["Ingress"]).To(BeEquivalentTo(3))
		Expect(spyMetrics.counters["Dropped"]).To(BeEquivalentTo(0))
		Expect(spyMetrics.counters["Egress"]).To(BeEquivalentTo(3))
	})

	It("increments dropped metric when write errors", func() {
		spyWriter := newSpyWriter()
		spyWriter.result = errors.New("write error")
		spyMetrics := newSpyMetrics()
		v := syslog.NewVisitor(spyWriter, spyMetrics)
		envs := []*loggregator_v2.Envelope{
			{SourceId: "source-1"},
			{SourceId: "source-2"},
			{SourceId: "source-3"},
		}

		v(envs)

		Expect(spyMetrics.counters["Ingress"]).To(BeEquivalentTo(3))
		Expect(spyMetrics.counters["Dropped"]).To(BeEquivalentTo(3))
		Expect(spyMetrics.counters["Egress"]).To(BeEquivalentTo(0))
	})
})

type spyWriter struct {
	called     int
	calledWith []*loggregator_v2.Envelope
	result     error
}

func newSpyWriter() *spyWriter {
	return &spyWriter{}
}

func (s *spyWriter) Write(e *loggregator_v2.Envelope) error {
	s.called++
	s.calledWith = append(s.calledWith, e)
	return s.result
}

type spyMetrics struct {
	counters map[string]uint64
}

func newSpyMetrics() *spyMetrics {
	return &spyMetrics{
		counters: make(map[string]uint64),
	}
}

func (s *spyMetrics) NewCounter(k string) func(uint64) {
	return func(i uint64) {
		s.counters[k] += i
	}
}

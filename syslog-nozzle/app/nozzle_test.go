package app_test

import (
	"bufio"
	"context"
	"fmt"
	"net"
	"runtime"
	"sync"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	loggregator "code.cloudfoundry.org/go-loggregator"
	"code.cloudfoundry.org/go-loggregator/rpc/loggregator_v2"
	"code.cloudfoundry.org/loggregator-tools/syslog-nozzle/app"
	"github.com/prometheus/client_golang/prometheus"
)

var _ = Describe("Nozzle", func() {
	It("sends logs from rlp to syslog", func() {
		spyStreamConnector := newSpyStreamConnector()
		syslogListener, err := net.Listen("tcp", ":0")
		Expect(err).ToNot(HaveOccurred())
		defer syslogListener.Close()
		syslogAddr := syslogListener.Addr().String()
		spyEgressedCounter := &spyCounter{}
		spyIgnoredCounter := &spyCounter{}
		nozzle := app.NewNozzle(
			spyStreamConnector,
			syslogAddr,
			"some-shard-id",
			app.WithEgressedCounter(spyEgressedCounter),
			app.WithIgnoredCounter(spyIgnoredCounter),
		)
		spyStreamConnector.addEnvelopeWithTags(0, "test-source-id-1", map[string]string{"namespace": "ns1"})
		spyStreamConnector.addEnvelopeWithTags(0, "test-source-id-2", map[string]string{"namespace": "ns2"})
		spyStreamConnector.addEnvelopeWithTags(0, "test-source-id-3", map[string]string{"namespace": "ns3"})

		go nozzle.Start()
		defer nozzle.Close()

		conn, err := syslogListener.Accept()
		Expect(err).ToNot(HaveOccurred())
		buf := bufio.NewReader(conn)

		errs := make(chan error, 100)
		msgs := make(chan string, 100)
		done := make(chan struct{})
		defer close(done)
		go func() {
			for {
				select {
				case <-done:
					return
				default:
					actual, err := buf.ReadString('\n')
					if err != nil {
						errs <- err
					}
					msgs <- actual
				}
			}

		}()
		Eventually(errs).Should(HaveLen(0))
		Eventually(msgs).Should(HaveLen(3))
		expected := fmt.Sprintf("85 <14>1 1970-01-01T00:00:00+00:00 - test-source-id-1 - - [tags@47450 namespace=\"ns1\"] \n")
		Expect(<-msgs).To(Equal(expected))
		expected = fmt.Sprintf("85 <14>1 1970-01-01T00:00:00+00:00 - test-source-id-2 - - [tags@47450 namespace=\"ns2\"] \n")
		Expect(<-msgs).To(Equal(expected))
		expected = fmt.Sprintf("85 <14>1 1970-01-01T00:00:00+00:00 - test-source-id-3 - - [tags@47450 namespace=\"ns3\"] \n")
		Expect(<-msgs).To(Equal(expected))

		Expect(spyEgressedCounter.read()).To(Equal(3))
		Expect(spyIgnoredCounter.read()).To(Equal(0))
	})

	It("sends logs for the specified namespace to syslog", func() {
		spyStreamConnector := newSpyStreamConnector()
		syslogListener, err := net.Listen("tcp", ":0")
		Expect(err).ToNot(HaveOccurred())
		defer syslogListener.Close()
		syslogAddr := syslogListener.Addr().String()
		spyEgressCounter := &spyCounter{}
		spyIgnoredCounter := &spyCounter{}
		nozzle := app.NewNozzle(
			spyStreamConnector,
			syslogAddr,
			"some-shard-id",
			app.WithEgressedCounter(spyEgressCounter),
			app.WithIgnoredCounter(spyIgnoredCounter),
			app.WithNamespace("ns1"),
		)

		spyStreamConnector.addEnvelopeWithTags(0, "ns1/rt1/rn1", map[string]string{"namespace": "ns1"})
		spyStreamConnector.addEnvelopeWithTags(0, "ns2/rt1/rn1", map[string]string{"namespace": "ns2"})
		spyStreamConnector.addEnvelopeWithTags(0, "ns1/rt2/rn2", map[string]string{"namespace": "ns1"})

		go nozzle.Start()
		defer nozzle.Close()

		conn, err := syslogListener.Accept()
		Expect(err).ToNot(HaveOccurred())
		buf := bufio.NewReader(conn)

		errs := make(chan error, 100)
		msgs := make(chan string, 100)
		done := make(chan struct{})
		defer close(done)
		go func() {
			for {
				select {
				case <-done:
					return
				default:
					actual, err := buf.ReadString('\n')
					if err != nil {
						errs <- err
					}
					msgs <- actual
				}
			}

		}()
		Eventually(msgs).Should(HaveLen(2))
		Eventually(errs).Should(HaveLen(0))
		expected := fmt.Sprintf("80 <14>1 1970-01-01T00:00:00+00:00 - ns1/rt1/rn1 - - [tags@47450 namespace=\"ns1\"] \n")
		Expect(<-msgs).To(Equal(expected))
		expected = fmt.Sprintf("80 <14>1 1970-01-01T00:00:00+00:00 - ns1/rt2/rn2 - - [tags@47450 namespace=\"ns1\"] \n")
		Expect(<-msgs).To(Equal(expected))

		Expect(spyEgressCounter.read()).To(Equal(2))
		Expect(spyIgnoredCounter.read()).To(Equal(1))
	})

	It("can stop", func() {
		spyStreamConnector := newSpyStreamConnector()
		syslogListener, err := net.Listen("tcp", ":0")
		Expect(err).ToNot(HaveOccurred())
		defer syslogListener.Close()
		syslogAddr := syslogListener.Addr().String()
		nozzle := app.NewNozzle(
			spyStreamConnector,
			syslogAddr,
			"some-shard-id",
		)
		start := runtime.NumGoroutine()

		go nozzle.Start()
		nozzle.Close()

		end := runtime.NumGoroutine()
		Expect(start).To(Equal(end))

		Eventually(spyStreamConnector.contexts()[0].Done()).Should(BeClosed())
	})

	It("drops messages that are not able to be encoded as syslog", func() {
		spyStreamConnector := newSpyStreamConnector()
		syslogListener, err := net.Listen("tcp", ":0")
		Expect(err).ToNot(HaveOccurred())
		defer syslogListener.Close()
		syslogAddr := syslogListener.Addr().String()
		spyCounter := &spyCounter{}
		nozzle := app.NewNozzle(
			spyStreamConnector,
			syslogAddr,
			"some-shard-id",
			app.WithConversionErrorCounter(spyCounter),
		)
		spyStreamConnector.addEnvelope(0, "test-source-id-1")
		spyStreamConnector.addBadEnvelope()
		spyStreamConnector.addEnvelope(0, "test-source-id-2")

		go nozzle.Start()
		defer nozzle.Close()

		conn, err := syslogListener.Accept()
		Expect(err).ToNot(HaveOccurred())
		buf := bufio.NewReader(conn)

		actual, err := buf.ReadString('\n')
		Expect(err).ToNot(HaveOccurred())
		expected := fmt.Sprintf("58 <14>1 1970-01-01T00:00:00+00:00 - test-source-id-1 - - - \n")
		Expect(actual).To(Equal(expected))

		actual, err = buf.ReadString('\n')
		Expect(err).ToNot(HaveOccurred())
		expected = fmt.Sprintf("58 <14>1 1970-01-01T00:00:00+00:00 - test-source-id-2 - - - \n")
		Expect(actual).To(Equal(expected))

		Expect(spyCounter.read()).To(Equal(1))
	})

	It("drops messages that are not able to be written to syslog", func() {
		spyStreamConnector := newSpyStreamConnector()
		syslogListener, err := net.Listen("tcp", ":0")
		Expect(err).ToNot(HaveOccurred())
		syslogAddr := syslogListener.Addr().String()
		spyEgressedCounter := &spyCounter{}
		spyDroppedCounter := &spyCounter{}
		nozzle := app.NewNozzle(
			spyStreamConnector,
			syslogAddr,
			"some-shard-id",
			app.WithEgressedCounter(spyEgressedCounter),
			app.WithDroppedCounter(spyDroppedCounter),
		)

		go nozzle.Start()
		defer nozzle.Close()

		spyStreamConnector.addEnvelope(0, "test-source-id-1")

		Eventually(spyEgressedCounter.read).Should(Equal(1))

		conn, err := syslogListener.Accept()
		Expect(err).ToNot(HaveOccurred())
		err = conn.Close()
		Expect(err).ToNot(HaveOccurred())
		err = syslogListener.Close()
		Expect(err).ToNot(HaveOccurred())

		spyStreamConnector.addEnvelope(0, "test-source-id-1")
		spyStreamConnector.addEnvelope(0, "test-source-id-2")
		spyStreamConnector.addEnvelope(0, "test-source-id-3")

		Eventually(spyDroppedCounter.read).Should(Equal(3))
	})

	It("returns error when unable to connect to syslog destination", func() {
		nozzle := app.NewNozzle(
			newSpyStreamConnector(),
			"unknown-syslog-addr",
			"some-shard-id",
		)

		Expect(nozzle.Start()).ToNot(Succeed())
	})
})

type spyStreamConnector struct {
	mu        sync.Mutex
	requests_ []*loggregator_v2.EgressBatchRequest
	contexts_ []context.Context
	envelopes chan []*loggregator_v2.Envelope
}

func newSpyStreamConnector() *spyStreamConnector {
	return &spyStreamConnector{
		envelopes: make(chan []*loggregator_v2.Envelope, 100),
	}
}

func (s *spyStreamConnector) Stream(ctx context.Context, req *loggregator_v2.EgressBatchRequest) loggregator.EnvelopeStream {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.requests_ = append(s.requests_, req)
	s.contexts_ = append(s.contexts_, ctx)

	return func() []*loggregator_v2.Envelope {
		select {
		case e := <-s.envelopes:
			return e
		default:
			return nil
		}
	}
}

func (s *spyStreamConnector) requests() []*loggregator_v2.EgressBatchRequest {
	s.mu.Lock()
	defer s.mu.Unlock()

	reqs := make([]*loggregator_v2.EgressBatchRequest, len(s.requests_))
	copy(reqs, s.requests_)

	return reqs
}

func (s *spyStreamConnector) contexts() []context.Context {
	s.mu.Lock()
	defer s.mu.Unlock()

	ctxs := make([]context.Context, len(s.contexts_))
	copy(ctxs, s.contexts_)

	return ctxs
}

func (c *spyStreamConnector) addEnvelope(timestamp int64, sourceID string) {
	c.envelopes <- []*loggregator_v2.Envelope{
		{
			Timestamp: timestamp,
			SourceId:  sourceID,
		},
	}
}

func (c *spyStreamConnector) addEnvelopeWithTags(timestamp int64, sourceID string, tags map[string]string) {
	c.envelopes <- []*loggregator_v2.Envelope{
		{
			Timestamp: timestamp,
			SourceId:  sourceID,
			Tags:      tags,
		},
	}
}

func (c *spyStreamConnector) addBadEnvelope() {
	c.envelopes <- []*loggregator_v2.Envelope{
		{
			SourceId: "\x01",
		},
	}
}

type spyCounter struct {
	prometheus.Counter
	mu    sync.Mutex
	count int
}

func (c *spyCounter) Inc() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.count++
}

func (c *spyCounter) read() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.count
}

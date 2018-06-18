package app_test

import (
	"bufio"
	"context"
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
	It("sends all logs from rlp to syslog", func() {
		spyStreamConnector, _, _, _, _ := setupSpies()
		spyDrain := newSpyDrain()
		defer spyDrain.stop()

		nozzle := app.NewNozzle(
			spyStreamConnector,
			[]app.Drain{
				{
					All: true,
					URL: spyDrain.url(),
				},
			},
			"",
		)
		spyStreamConnector.addEnvelope("test-source-id-1", map[string]string{"namespace": "ns1"})
		spyStreamConnector.addEnvelope("test-source-id-2", map[string]string{"namespace": "ns2"})
		spyStreamConnector.addEnvelope("test-source-id-3", map[string]string{"namespace": "ns3"})

		go nozzle.Start()
		defer nozzle.Close()

		spyDrain.expectReceived(
			`85 <14>1 1970-01-01T00:00:00+00:00 - test-source-id-1 - - [tags@47450 namespace="ns1"] `+"\n",
			`85 <14>1 1970-01-01T00:00:00+00:00 - test-source-id-2 - - [tags@47450 namespace="ns2"] `+"\n",
			`85 <14>1 1970-01-01T00:00:00+00:00 - test-source-id-3 - - [tags@47450 namespace="ns3"] `+"\n",
		)
	})

	It("routes logs to the specified drains based on namespace tag", func() {
		spyStreamConnector, _, _, _, _ := setupSpies()
		spyDrain1 := newSpyDrain()
		defer spyDrain1.stop()
		spyDrain2 := newSpyDrain()
		defer spyDrain2.stop()

		nozzle := app.NewNozzle(
			spyStreamConnector,
			[]app.Drain{
				{
					Namespace: "ns1",
					URL:       spyDrain1.url(),
				},
				{
					Namespace: "ns2",
					URL:       spyDrain2.url(),
				},
			},
			"",
		)

		spyStreamConnector.addEnvelope("ns1/rt1/rn1")
		spyStreamConnector.addEnvelope("ns1/rt1/rn1", map[string]string{"namespace": "ns1"})
		spyStreamConnector.addEnvelope("ns2/rt1/rn1", map[string]string{"namespace": "ns2"})
		spyStreamConnector.addEnvelope("ns1/rt2/rn2", map[string]string{"namespace": "ns1"})

		go nozzle.Start()
		defer nozzle.Close()

		spyDrain1.expectReceived(
			`80 <14>1 1970-01-01T00:00:00+00:00 - ns1/rt1/rn1 - - [tags@47450 namespace="ns1"] `+"\n",
			`80 <14>1 1970-01-01T00:00:00+00:00 - ns1/rt2/rn2 - - [tags@47450 namespace="ns1"] `+"\n",
		)
		spyDrain2.expectReceived(
			`80 <14>1 1970-01-01T00:00:00+00:00 - ns2/rt1/rn1 - - [tags@47450 namespace="ns2"] ` + "\n",
		)
	})

	It("can stop", func() {
		spyStreamConnector, _, _, _, _ := setupSpies()
		spyDrain := newSpyDrain()
		defer spyDrain.stop()

		nozzle := app.NewNozzle(
			spyStreamConnector,
			[]app.Drain{
				{
					All: true,
					URL: spyDrain.url(),
				},
			},
			"",
		)
		start := runtime.NumGoroutine()

		go nozzle.Start()

		err := nozzle.Close()
		Expect(err).ToNot(HaveOccurred())
		end := runtime.NumGoroutine()
		Expect(start).To(Equal(end))
		Eventually(spyStreamConnector.contexts()[0].Done()).Should(BeClosed())
	})

	It("drops messages that are not able to be encoded as syslog", func() {
		spyStreamConnector, _, _, _, spyConversionErrorCounter := setupSpies()
		spyDrain := newSpyDrain()
		defer spyDrain.stop()

		nozzle := app.NewNozzle(
			spyStreamConnector,
			[]app.Drain{
				{
					All: true,
					URL: spyDrain.url(),
				},
			},
			"",
			app.WithConversionErrorCounter(spyConversionErrorCounter),
		)
		spyStreamConnector.addEnvelope("test-source-id-1")
		spyStreamConnector.addEnvelope("\x01")
		spyStreamConnector.addEnvelope("test-source-id-2")

		go nozzle.Start()
		defer nozzle.Close()

		spyDrain.expectReceived(
			"58 <14>1 1970-01-01T00:00:00+00:00 - test-source-id-1 - - - \n",
			"58 <14>1 1970-01-01T00:00:00+00:00 - test-source-id-2 - - - \n",
		)
	})

	Describe("metrics", func() {
		It("increments egressed counter when message is written", func() {
			spyStreamConnector, spyEgressedCounter, _, _, _ := setupSpies()
			spyDrain := newSpyDrain()
			defer spyDrain.stop()

			nozzle := app.NewNozzle(
				spyStreamConnector,
				[]app.Drain{
					{
						All: true,
						URL: spyDrain.url(),
					},
				},
				"",
				app.WithEgressedCounter(spyEgressedCounter),
			)
			spyStreamConnector.addEnvelope("test-source-id-1", map[string]string{"namespace": "ns1"})

			go nozzle.Start()
			defer nozzle.Close()

			Eventually(spyEgressedCounter.read).Should(Equal(1))
		})

		It("increments ignored counter for envelopes that do not match desired namespace", func() {
			spyStreamConnector, _, spyIgnoredCounter, _, _ := setupSpies()
			spyDrain := newSpyDrain()
			defer spyDrain.stop()

			nozzle := app.NewNozzle(
				spyStreamConnector,
				[]app.Drain{
					{
						Namespace: "ns1",
						URL:       spyDrain.url(),
					},
				},
				"",
				app.WithIgnoredCounter(spyIgnoredCounter),
			)

			spyStreamConnector.addEnvelope("ns2/rt1/rn1", map[string]string{"namespace": "ns2"})

			go nozzle.Start()
			defer nozzle.Close()

			Eventually(spyIgnoredCounter.read).Should(Equal(1))
		})

		It("increments ignored counter for envelopes that do not have a namespace tag", func() {
			spyStreamConnector, _, spyIgnoredCounter, _, _ := setupSpies()
			spyDrain := newSpyDrain()
			defer spyDrain.stop()

			nozzle := app.NewNozzle(
				spyStreamConnector,
				[]app.Drain{
					{
						All: true,
						URL: spyDrain.url(),
					},
				},
				"",
				app.WithIgnoredCounter(spyIgnoredCounter),
			)

			spyStreamConnector.addEnvelope("test-source-id-1")

			go nozzle.Start()
			defer nozzle.Close()

			Eventually(spyIgnoredCounter.read).Should(Equal(1))
		})

		It("increments dropped counter when unable to write to syslog drain", func() {
			spyStreamConnector, spyEgressedCounter, spyDroppedCounter, _, _ := setupSpies()
			spyDrain := newSpyDrain()
			defer spyDrain.stop()

			nozzle := app.NewNozzle(
				spyStreamConnector,
				[]app.Drain{
					{
						All: true,
						URL: spyDrain.url(),
					},
				},
				"",
				app.WithEgressedCounter(spyEgressedCounter),
				app.WithDroppedCounter(spyDroppedCounter),
			)

			go nozzle.Start()
			defer nozzle.Close()

			spyStreamConnector.addEnvelope("test-source-id-1")
			Eventually(spyEgressedCounter.read).Should(Equal(1))

			conn := spyDrain.accept()
			err := conn.Close()
			Expect(err).ToNot(HaveOccurred())
			spyDrain.stop()

			spyStreamConnector.addEnvelope("test-source-id-1")
			spyStreamConnector.addEnvelope("test-source-id-2")
			spyStreamConnector.addEnvelope("test-source-id-3")

			Eventually(spyDroppedCounter.read).Should(Equal(3))
		})

		It("increments conversion error counter when envelope cannot be converted to syslog", func() {
			spyStreamConnector, _, _, _, spyConversionErrorCounter := setupSpies()
			spyDrain := newSpyDrain()
			defer spyDrain.stop()

			nozzle := app.NewNozzle(
				spyStreamConnector,
				[]app.Drain{
					{
						All: true,
						URL: spyDrain.url(),
					},
				},
				"",
				app.WithConversionErrorCounter(spyConversionErrorCounter),
			)
			spyStreamConnector.addEnvelope("\x01")

			go nozzle.Start()
			defer nozzle.Close()

			Eventually(spyConversionErrorCounter.read).Should(Equal(1))
		})
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

func (c *spyStreamConnector) addEnvelope(sourceID string, tags ...map[string]string) {
	e := &loggregator_v2.Envelope{
		SourceId: sourceID,
	}
	if len(tags) > 0 {
		e.Tags = tags[0]
	}
	c.envelopes <- []*loggregator_v2.Envelope{e}
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

func setupSpies() (
	spyStreamConnector *spyStreamConnector,
	spyEgressedCounter *spyCounter,
	spyIgnoredCounter *spyCounter,
	spyDroppedCounter *spyCounter,
	spyConversionErrorCounter *spyCounter,
) {
	spyStreamConnector = newSpyStreamConnector()
	spyEgressedCounter = &spyCounter{}
	spyIgnoredCounter = &spyCounter{}
	spyDroppedCounter = &spyCounter{}
	spyConversionErrorCounter = &spyCounter{}
	return
}

type spyDrain struct {
	lis net.Listener
}

func newSpyDrain() *spyDrain {
	lis, err := net.Listen("tcp", ":0")
	Expect(err).ToNot(HaveOccurred())

	return &spyDrain{
		lis: lis,
	}
}

func (s *spyDrain) url() string {
	return s.lis.Addr().String()
}

func (s *spyDrain) stop() {
	s.lis.Close()
}

func (s *spyDrain) accept() net.Conn {
	conn, err := s.lis.Accept()
	Expect(err).ToNot(HaveOccurred())
	return conn
}

func (s *spyDrain) expectReceived(msgs ...string) {
	conn := s.accept()
	defer conn.Close()
	buf := bufio.NewReader(conn)

	for _, expected := range msgs {
		actual, err := buf.ReadString('\n')
		Expect(err).ToNot(HaveOccurred())
		Expect(actual).To(Equal(expected))
	}
}

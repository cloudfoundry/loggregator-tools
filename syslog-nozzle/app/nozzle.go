package app

import (
	"context"
	"fmt"
	"io"
	"net"

	loggregator "code.cloudfoundry.org/go-loggregator"
	"code.cloudfoundry.org/go-loggregator/rpc/loggregator_v2"
	"github.com/prometheus/client_golang/prometheus"
)

type noopCounter struct {
	prometheus.Counter
}

func (noopCounter) Inc() {}

type Nozzle struct {
	sc                 StreamConnector
	drains             []Drain
	drainWriters       map[string]io.Writer
	globalDrainWriters []io.Writer
	shardID            string

	stop chan struct{}
	done chan struct{}

	egressedCounter   prometheus.Counter
	droppedCounter    prometheus.Counter
	ignoredEnvCounter prometheus.Counter
	envConvertCounter prometheus.Counter
}

type StreamConnector interface {
	Stream(context.Context, *loggregator_v2.EgressBatchRequest) loggregator.EnvelopeStream
}

type NozzleOption func(*Nozzle)

func WithEgressedCounter(c prometheus.Counter) NozzleOption {
	return func(n *Nozzle) {
		n.egressedCounter = c
	}
}

func WithDroppedCounter(c prometheus.Counter) NozzleOption {
	return func(n *Nozzle) {
		n.droppedCounter = c
	}
}

func WithIgnoredCounter(c prometheus.Counter) NozzleOption {
	return func(n *Nozzle) {
		n.ignoredEnvCounter = c
	}
}

func WithConversionErrorCounter(c prometheus.Counter) NozzleOption {
	return func(n *Nozzle) {
		n.envConvertCounter = c
	}
}

func NewNozzle(
	sc StreamConnector,
	drains []Drain,
	shardID string,
	opts ...NozzleOption,
) *Nozzle {
	n := &Nozzle{
		sc:           sc,
		drains:       drains,
		drainWriters: make(map[string]io.Writer),
		shardID:      shardID,
		stop:         make(chan struct{}),
		done:         make(chan struct{}),

		egressedCounter:   noopCounter{},
		droppedCounter:    noopCounter{},
		ignoredEnvCounter: noopCounter{},
		envConvertCounter: noopCounter{},
	}

	for _, o := range opts {
		o(n)
	}

	return n
}

func (n *Nozzle) Start() error {
	defer close(n.done)

	for _, d := range n.drains {
		conn, err := net.Dial("tcp", d.URL)
		if err != nil {
			// TODO: handle conn errors, reconnect?
			// TODO: Write test that force "continue"
			// continue
		}
		defer conn.Close()
		if d.All {
			n.globalDrainWriters = append(n.globalDrainWriters, conn)
			continue
		}
		n.drainWriters[d.Namespace] = conn
	}

	req := &loggregator_v2.EgressBatchRequest{
		ShardId: n.shardID,
		Selectors: []*loggregator_v2.Selector{
			{
				Message: &loggregator_v2.Selector_Log{
					Log: &loggregator_v2.LogSelector{},
				},
			},
			{
				Message: &loggregator_v2.Selector_Counter{
					Counter: &loggregator_v2.CounterSelector{},
				},
			},
			{
				Message: &loggregator_v2.Selector_Gauge{
					Gauge: &loggregator_v2.GaugeSelector{},
				},
			},
			{
				Message: &loggregator_v2.Selector_Timer{
					Timer: &loggregator_v2.TimerSelector{},
				},
			},
			{
				Message: &loggregator_v2.Selector_Event{
					Event: &loggregator_v2.EventSelector{},
				},
			},
		},
		UsePreferredTags: true,
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	stream := n.sc.Stream(ctx, req)

	for {
		for _, e := range stream() {
			msgs, err := e.Syslog()
			if err != nil {
				n.envConvertCounter.Inc()
			}

			for _, w := range n.globalDrainWriters {
				n.write(w, msgs)
			}

			w, ok := n.drainWriters[namespace(e)]
			if !ok {
				n.ignoredEnvCounter.Inc()
				continue
			}
			n.write(w, msgs)
		}

		select {
		case <-n.stop:
			// TODO: Write test that force "return nil" vs "return err"
			return nil
		default:
		}
	}
}

func (n *Nozzle) write(w io.Writer, msgs [][]byte) {
	for _, m := range msgs {
		// TODO: test io timeout to drive out conn.SetWriteDeadline()
		// TODO: what happens if we write a partial message, we
		//       probably shouldn't reuse the conn
		_, err := fmt.Fprintf(w, "%d %s", len(m), m)
		if err != nil {
			n.droppedCounter.Inc()
			// TODO: drive out
			// continue
		}
		n.egressedCounter.Inc()
	}
}

func (n *Nozzle) Close() error {
	close(n.stop)
	<-n.done
	return nil
}

func namespace(e *loggregator_v2.Envelope) string {
	tags := e.GetTags()
	if tags == nil {
		return ""
	}
	return tags["namespace"]
}

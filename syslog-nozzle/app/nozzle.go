package app

import (
	"context"
	"fmt"
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
	sc        StreamConnector
	addr      string
	shardID   string
	namespace string

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

func WithNamespace(ns string) NozzleOption {
	return func(n *Nozzle) {
		n.namespace = ns
	}
}

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

func WithIgnoredEnvelopeCounter(c prometheus.Counter) NozzleOption {
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
	syslogAddr, shardID string,
	opts ...NozzleOption,
) *Nozzle {
	n := &Nozzle{
		sc:      sc,
		addr:    syslogAddr,
		shardID: shardID,
		stop:    make(chan struct{}),
		done:    make(chan struct{}),

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

	conn, err := net.Dial("tcp", n.addr)
	if err != nil {
		return err
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
			if !n.fromNamespace(e) {
				n.ignoredEnvCounter.Inc()
				continue
			}
			slm, err := e.Syslog()
			if err != nil {
				n.envConvertCounter.Inc()
			}
			for _, m := range slm {
				// TODO: test io timeout to drive out conn.SetWriteDeadline()
				// TODO: what happens if we write a partial message, we
				//       probably shouldn't reuse the conn
				_, err := fmt.Fprintf(conn, "%d %s", len(m), m)
				if err != nil {
					n.droppedCounter.Inc()
					continue
				}
				n.egressedCounter.Inc()
			}
		}

		select {
		case <-n.stop:
			return nil
		default:
		}
	}
}

func (n *Nozzle) Close() error {
	close(n.stop)
	<-n.done
	return nil
}

func (n *Nozzle) fromNamespace(e *loggregator_v2.Envelope) bool {
	if n.namespace == "" {
		return true
	}
	ns := e.GetTags()["namespace"]
	return ns == n.namespace
}

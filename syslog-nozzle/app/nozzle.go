package app

import (
	"context"
	"fmt"
	"net"
	"time"

	loggregator "code.cloudfoundry.org/go-loggregator"
	"code.cloudfoundry.org/go-loggregator/rpc/loggregator_v2"
	"github.com/prometheus/client_golang/prometheus"
)

type noopCounter struct {
	prometheus.Counter
}

func (noopCounter) Inc() {}

type Nozzle struct {
	sc               StreamConnector
	drains           []Drain
	drainConns       map[string]net.Conn
	globalDrainConns []net.Conn
	shardID          string

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
		sc:         sc,
		drains:     drains,
		drainConns: make(map[string]net.Conn),
		shardID:    shardID,
		stop:       make(chan struct{}),
		done:       make(chan struct{}),

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
			// TODO: We should trigger reconnect here vs ignoring the drain.
			continue
		}
		defer conn.Close()
		if d.All {
			n.globalDrainConns = append(n.globalDrainConns, conn)
			continue
		}
		n.drainConns[d.Namespace] = conn
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

			for _, c := range n.globalDrainConns {
				n.write(c, msgs)
			}

			c, ok := n.drainConns[e.GetTags()["namespace"]]
			if !ok {
				n.ignoredEnvCounter.Inc()
				continue
			}
			n.write(c, msgs)
		}

		select {
		case <-n.stop:
			return nil
		default:
		}
	}
}

func (n *Nozzle) write(conn net.Conn, msgs [][]byte) {
	for _, m := range msgs {
		conn.SetWriteDeadline(time.Now().Add(5 * time.Millisecond))
		_, err := fmt.Fprintf(conn, "%d %s", len(m), m)
		if err != nil {
			// TODO: When we get an error back from writing to conn, we should
			// probably reconnect.
			n.droppedCounter.Inc()
			continue
		}
		n.egressedCounter.Inc()
	}
}

func (n *Nozzle) Close() error {
	close(n.stop)
	<-n.done
	return nil
}

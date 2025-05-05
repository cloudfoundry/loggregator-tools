// envelopeemitter: a tool to emit envelopes via v2 gRPC
package main

import (
	"context"
	"flag"
	"log"
	"time"

	"code.cloudfoundry.org/go-loggregator/v10/rpc/loggregator_v2"
	"code.cloudfoundry.org/loggregator-release/src/plumbing"

	"google.golang.org/grpc"
)

var (
	addr     = flag.String("addr", "localhost:3458", "address to connect for gRPC")
	certFile = flag.String("cert", "", "cert to use to connect for gRPC")
	keyFile  = flag.String("key", "", "key to use to connect for gRPC")
	caFile   = flag.String("ca", "", "ca cert to use to connect for gRPC")
	srcID    = flag.String("source-id", "", "envelope source ID")
	value    = flag.Float64("value", 0.0, "a value")
)

func main() {
	flag.Parse()
	creds, err := plumbing.NewClientCredentials(*certFile, *keyFile, *caFile, "metron")
	if err != nil {
		log.Fatal(err)
	}

	conn, err := grpc.NewClient(*addr, grpc.WithTransportCredentials(creds))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close() //nolint:errcheck
	c := loggregator_v2.NewIngressClient(conn)

	// create env
	env := &loggregator_v2.Envelope{
		SourceId: *srcID,
		Message: &loggregator_v2.Envelope_Gauge{
			Gauge: &loggregator_v2.Gauge{
				Metrics: map[string]*loggregator_v2.GaugeValue{
					"test-metric": &loggregator_v2.GaugeValue{
						Value: *value,
					},
				},
			},
		},
	}
	sender, err := c.Sender(context.TODO())
	if err != nil {
		log.Fatal(err)
	}

	for {
		err := sender.Send(env)
		if err != nil {
			log.Fatal(err)
		}
		time.Sleep(time.Second)
		log.Printf("emmiting a gauge")
	}
}

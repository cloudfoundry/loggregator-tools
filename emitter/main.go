// emitter: a tool to emit envelopes to the agent via v2 gRPC.
//
// see: https://hub.docker.com/r/loggregator/emitter/
//
package main

import (
	"context"
	"log"
	"time"

	"code.cloudfoundry.org/go-envstruct"
	"code.cloudfoundry.org/loggregator/plumbing"
	"code.cloudfoundry.org/loggregator/plumbing/v2"

	"google.golang.org/grpc"
)

type config struct {
	Target   string        `env:"TARGET,required"`
	CertFile string        `env:"CERT,required"`
	KeyFile  string        `env:"KEY,required"`
	CAFile   string        `env:"CA,required"`
	Rate     time.Duration `env:"RATE,required"`
}

func main() {
	c := config{}
	err := envstruct.Load(&c)
	if err != nil {
		log.Fatal(err)
	}

	creds, err := plumbing.NewClientCredentials(
		c.CertFile,
		c.KeyFile,
		c.CAFile,
		"agent",
	)
	if err != nil {
		log.Fatal(err)
	}

	conn, err := grpc.Dial(c.Target, grpc.WithTransportCredentials(creds))
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		err := conn.Close()
		if err != nil {
			log.Fatal()
		}
	}()
	client := loggregator_v2.NewIngressClient(conn)

	env := &loggregator_v2.Envelope{
		Message: &loggregator_v2.Envelope_Counter{
			Counter: &loggregator_v2.Counter{
				Name:  "some-counter",
				Delta: 5,
			},
		},
	}
	sender, err := client.Sender(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	for {
		log.Printf("emmiting a counter")
		err := sender.Send(env)
		if err != nil {
			log.Fatal(err)
		}
		time.Sleep(c.Rate)
	}
}

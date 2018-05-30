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
	"code.cloudfoundry.org/go-loggregator"
	"code.cloudfoundry.org/go-loggregator/rpc/loggregator_v2"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type config struct {
	Target   string        `env:"TARGET,required"`
	CertFile string        `env:"CERT,required"`
	KeyFile  string        `env:"KEY,required"`
	CAFile   string        `env:"CA,required"`
	Rate     time.Duration `env:"RATE,required"`
	SourceID string        `env:"SOURCE_ID"`
}

func main() {
	c := config{}
	err := envstruct.Load(&c)
	if err != nil {
		log.Fatal(err)
	}

	creds, err := loggregator.NewIngressTLSConfig(
		c.CAFile,
		c.CertFile,
		c.KeyFile,
	)
	if err != nil {
		log.Fatal(err)
	}

	conn, err := grpc.Dial(c.Target, grpc.WithTransportCredentials(
		credentials.NewTLS(creds),
	))
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
		SourceId: c.SourceID,
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
		env.Timestamp = time.Now().UnixNano()
		err := sender.Send(env)
		if err != nil {
			log.Fatal(err)
		}
		time.Sleep(c.Rate)
	}
}

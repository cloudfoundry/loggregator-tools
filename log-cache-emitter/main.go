package main

import (
	"context"
	"log"
	_ "net/http/pprof"
	"time"

	envstruct "code.cloudfoundry.org/go-envstruct"
	"code.cloudfoundry.org/go-log-cache/rpc/logcache_v1"
	"code.cloudfoundry.org/go-loggregator/rpc/loggregator_v2"
	"google.golang.org/grpc"
)

func main() {
	log.Print("Starting LogCache Nozzle...")
	defer log.Print("Closing LogCache Nozzle.")

	cfg, err := LoadConfig()
	if err != nil {
		log.Fatalf("invalid configuration: %s", err)
	}

	envstruct.WriteReport(cfg)

	conn, err := grpc.Dial(cfg.LogCacheAddr, grpc.WithTransportCredentials(
		cfg.LogCacheTLS.Credentials("log-cache"),
	))
	if err != nil {
		log.Fatalf("failed to dial %s: %s", cfg.LogCacheAddr, err)
	}
	client := logcache_v1.NewIngressClient(conn)

	for {
		batch := []*loggregator_v2.Envelope{
			{
				Timestamp:  time.Now().UnixNano(),
				SourceId:   "jasonk",
				InstanceId: "jason",
				Tags: map[string]string{
					"foo": "bar",
				},
				Message: &loggregator_v2.Envelope_Log{
					Log: &loggregator_v2.Log{
						Payload: []byte("warren is cool"),
						Type:    loggregator_v2.Log_ERR,
					},
				},
			},
		}
		log.Print("writing a batch")
		ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
		_, err := client.Send(ctx, &logcache_v1.SendRequest{
			Envelopes: &loggregator_v2.EnvelopeBatch{
				Batch: batch,
			},
		})

		if err != nil {
			log.Printf("failed to write envelopes: %s", err)
			continue
		}
		time.Sleep(time.Second)
	}
}

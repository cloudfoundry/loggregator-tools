package main

import (
	"context"
	"log"
	"sync"
	"time"

	"net/http"

	envstruct "code.cloudfoundry.org/go-envstruct"
	"code.cloudfoundry.org/go-log-cache/rpc/logcache_v1"
	"code.cloudfoundry.org/go-loggregator/rpc/loggregator_v2"
	"google.golang.org/grpc"
)

func main() {
	log.Print("Starting LogCache Emitter...")
	defer log.Print("Closing LogCache Emitter.")

	cfg, err := LoadConfig()
	if err != nil {
		log.Fatalf("invalid configuration: %s", err)
	}

	envstruct.WriteReport(cfg)

	http.HandleFunc("/emit", emitHandler(cfg))

	log.Printf("Listening on: %s", cfg.Addr)
	http.ListenAndServe(cfg.Addr, nil)
}

func emitHandler(cfg *Config) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		sourceIDs, ok := q["sourceIDs"]
		if !ok {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("sourceIDs are required\n"))
			return
		}

		emitLogs(cfg, sourceIDs)
	}
}

func emitLogs(cfg *Config, sourceIDs []string) {
	wg := &sync.WaitGroup{}
	for _, s := range sourceIDs {
		wg.Add(1)
		go sendLogs(cfg, wg, s)
	}
	wg.Wait()
}

func sendLogs(cfg *Config, wg *sync.WaitGroup, sourceID string) {
	conn, err := grpc.Dial(cfg.LogCacheAddr, grpc.WithTransportCredentials(
		cfg.LogCacheTLS.Credentials("log-cache"),
	))
	if err != nil {
		log.Fatalf("failed to dial %s: %s", cfg.LogCacheAddr, err)
	}

	client := logcache_v1.NewIngressClient(conn)
	log.Printf("Emitting Logs for %s", sourceID)
	for i := 0; i < 10000; i++ {
		batch := []*loggregator_v2.Envelope{
			{
				Timestamp: time.Now().UnixNano(),
				SourceId:  sourceID,
			},
		}
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
		time.Sleep(time.Millisecond)
	}

	wg.Done()
	log.Print("Done")
}

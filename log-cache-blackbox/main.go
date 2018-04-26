package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"sync/atomic"
	"time"

	envstruct "code.cloudfoundry.org/go-envstruct"
)

func main() {
	var cfg Config
	err := envstruct.Load(&cfg)
	if err != nil {
		log.Fatalf("failed to load config: %s", err)
	}

	err = envstruct.WriteReport(&cfg)
	if err != nil {
		log.Fatalf("failed to write config report: %s", err)
	}

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", cfg.Port), handler(cfg)))
}

func handler(cfg Config) http.Handler {
	var testRunning int64

	httpClient := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: cfg.SkipSSLValidation,
			},
		},
		Timeout: 5 * time.Second,
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		if !atomic.CompareAndSwapInt64(&testRunning, 0, 1) {
			w.WriteHeader(http.StatusLocked)
			return
		}
		defer atomic.StoreInt64(&testRunning, 0)

		switch r.URL.Path {
		case "/":
			latencyHandler(cfg, httpClient).ServeHTTP(w, r)
		case "/reliability":
			reliabilityHandler(cfg, httpClient).ServeHTTP(w, r)
		case "/group":
			groupLatencyHandler(cfg, httpClient).ServeHTTP(w, r)
		case "/group/reliability":
			groupReliabilityHandler(cfg, httpClient).ServeHTTP(w, r)
		default:
			w.WriteHeader(http.StatusNotFound)
			return
		}
	})
}

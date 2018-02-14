package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"time"

	envstruct "code.cloudfoundry.org/go-envstruct"
	logcache "code.cloudfoundry.org/go-log-cache"
	"code.cloudfoundry.org/loggregator-tools/log-cache-siege/internal/handlers"
)

func main() {
	log.Println("Starting Log Cache Siege...")
	defer log.Println("Closing Log Cache Siege.")

	cfg := loadConfig()

	httpClient := httpClient(cfg)
	addr := fmt.Sprintf(":%d", cfg.Port)
	logCacheClient := logcache.NewClient(
		cfg.LogCacheAddr,
		logcache.WithHTTPClient(httpClient))

	siegeHandler := handlers.NewSiege(addr, httpClient, logCacheClient)
	log.Fatal(http.ListenAndServe(addr, siegeHandler))
}

func httpClient(cfg config) *http.Client {
	return &http.Client{
		Timeout: 5 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: cfg.SkipSSLValidation,
			},
		},
	}
}

type config struct {
	Port              int    `env:"PORT,           required"`
	LogCacheAddr      string `env:"LOG_CACHE_ADDR, required"`
	SkipSSLValidation bool   `env:"SKIP_SSL_VALIDATION"`
}

func loadConfig() config {
	c := config{
		SkipSSLValidation: false,
	}

	if err := envstruct.Load(&c); err != nil {
		log.Fatal(err)
	}

	return c
}

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
	oauthHttpClient := logcache.NewOauth2HTTPClient(
		cfg.UAAAddr,
		cfg.UAAClient,
		cfg.UAAClientSecret,
		logcache.WithOauth2HTTPClient(httpClient),
	)

	logCacheClient := logcache.NewClient(
		cfg.LogCacheAddr,
		logcache.WithHTTPClient(oauthHttpClient),
	)

	siegeHandler := handlers.NewSiege(
		cfg.RequestSpinnerAddr,
		cfg.ConcurrentRequests,
		httpClient,
		logCacheClient,
	)

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", cfg.Port), siegeHandler))
}

func httpClient(cfg config) *http.Client {
	return &http.Client{
		Timeout: 1 * time.Minute,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: cfg.SkipSSLValidation,
			},
		},
	}
}

type config struct {
	Port               int    `env:"PORT,                 required"`
	LogCacheAddr       string `env:"LOG_CACHE_ADDR,       required"`
	RequestSpinnerAddr string `env:"REQUEST_SPINNER_ADDR, required"`
	ConcurrentRequests int    `env:"CONCURRENT_REQUESTS"`

	UAAAddr         string `env:"UAA_ADDR,          required"`
	UAAClient       string `env:"UAA_CLIENT,        required"`
	UAAClientSecret string `env:"UAA_CLIENT_SECRET, required"`

	SkipSSLValidation bool `env:"SKIP_SSL_VALIDATION"`
}

func loadConfig() config {
	c := config{
		SkipSSLValidation:  false,
		ConcurrentRequests: 25,
	}

	if err := envstruct.Load(&c); err != nil {
		log.Fatal(err)
	}

	return c
}

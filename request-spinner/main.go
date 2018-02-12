package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"time"

	envstruct "code.cloudfoundry.org/go-envstruct"
	logcache "code.cloudfoundry.org/go-log-cache"
	"code.cloudfoundry.org/loggregator-tools/request-spinner/internal/handlers"
)

func main() {
	log.Println("Starting request spinner...")
	defer log.Println("Closing request spinner.")

	cfg := loadConfig()

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: cfg.SkipSSLValidation},
	}

	client := logcache.NewClient(cfg.LogCacheAddr,
		logcache.WithHTTPClient(logcache.NewOauth2HTTPClient(
			cfg.UAAAddr,
			cfg.UAAClient,
			cfg.UAAClientSecret,
			logcache.WithOauth2HTTPClient(
				&http.Client{
					Transport: tr,
					Timeout:   5 * time.Second,
				},
			))),
	)
	http.ListenAndServe(
		fmt.Sprintf(":%d", cfg.Port),
		handlers.NewStarter(client),
	)

}

type config struct {
	Port         int    `env:"PORT,required"`
	LogCacheAddr string `env:"LOG_CACHE_ADDR,required"`

	UAAAddr         string `env:"UAA_ADDR,required"`
	UAAClient       string `env:"UAA_CLIENT,required"`
	UAAClientSecret string `env:"UAA_CLIENT_SECRET,required"`

	SkipSSLValidation bool `env:"SKIP_SSL_VALIDATION"`
}

func loadConfig() config {
	var c config

	if err := envstruct.Load(&c); err != nil {
		log.Fatal(err)
	}

	return c
}

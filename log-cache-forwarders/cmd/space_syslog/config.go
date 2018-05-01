package main

import (
	"log"
	"net/url"
	"time"

	envstruct "code.cloudfoundry.org/go-envstruct"
)

type Config struct {
	SpaceGUID      string `env:"SPACE_GUID",      required`
	SourceHostname string `env:"SOURCE_HOST_NAME, required"`

	LogCacheGroupName string `env:"GROUP_NAME",      required`

	CAPIURL string `env:"CAPI_URL", required`

	DrainName string `env:"DRAIN_NAME", required`
	DrainType string `env:"DRAIN_TYPE"`

	UAAURL       string `env:"UAA_URL,       required"`
	ClientID     string `env:"CLIENT_ID,     required"`
	ClientSecret string `env:"CLIENT_SECRET,          noreport"`

	Username string `env:"USERNAME, required"`
	Password string `env:"PASSWORD, required, noreport"`

	LogCacheHTTPAddr string   `env:"LOG_CACHE_HTTP_ADDR, required"`
	SyslogURL        *url.URL `env:"SYSLOG_URL,          required"`

	SkipCertVerify bool `env:"SKIP_CERT_VERIFY"`

	DialTimeout time.Duration `env:"DIAL_TIMEOUT"`
	IOTimeout   time.Duration `env:"IO_TIMEOUT"`
	KeepAlive   time.Duration `env:"KEEP_ALIVE"`
}

func LoadConfig() Config {
	cfg := Config{
		SkipCertVerify: false,
		KeepAlive:      10 * time.Second,
		DialTimeout:    5 * time.Second,
		IOTimeout:      time.Minute,
		DrainType:      "all",
	}
	if err := envstruct.Load(&cfg); err != nil {
		log.Fatalf("failed to load config from environment: %s", err)
	}

	return cfg
}

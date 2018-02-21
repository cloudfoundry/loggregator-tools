package main

import (
	"log"

	envstruct "code.cloudfoundry.org/go-envstruct"
)

type Config struct {
	// SourceIDBlacklist is a regular expression for envelopes to be excluded.
	SourceIDBlacklist string `env:"SOURCE_ID_BLACKLIST"`

	DatadogAPIKey string `env:"DATADOG_API_KEY, required, noreport"`
	MetricHost    string `env:"METRIC_HOST, required"`

	// DatadogTags are a comma separated list of tags to be set on each
	// metric.
	DatadogTags string `env:"DATADOG_TAGS"`

	UAAAddr      string `env:"UAA_ADDR,        required"`
	ClientID     string `env:"CLIENT_ID,       required"`
	ClientSecret string `env:"CLIENT_SECRET,   required, noreport"`

	LogCacheHTTPAddr  string `env:"LOG_CACHE_HTTP_ADDR,  required"`
	LogCacheGroupName string `env:"LOG_CACHE_GROUP_NAME, required"`

	SkipCertVerify bool `env:"SKIP_CERT_VERIFY"`
}

func LoadConfig() Config {
	cfg := Config{
		SkipCertVerify: false,
	}
	if err := envstruct.Load(&cfg); err != nil {
		log.Fatalf("failed to load config from environment: %s", err)
	}

	return cfg
}

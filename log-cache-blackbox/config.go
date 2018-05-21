package main

import (
	"encoding/json"
	"net/url"
	"time"
)

type Config struct {
	LogCacheURL *url.URL `env:"LOG_CACHE_URL, required, report"`

	Port     int     `env:"PORT,             report, required"`
	SourceID string  `env:"SOURCE_ID,        report"`
	VCapApp  VCapApp `env:"VCAP_APPLICATION, report"`

	UAAAddr         string `env:"UAA_ADDR,        report"`
	UAAClient       string `env:"UAA_CLIENT       "`
	UAAClientSecret string `env:"UAA_CLIENT_SECRET"`

	SkipSSLValidation bool `env:"SKIP_SSL_VALIDATION, report"`

	WalkDelay time.Duration `env:WALK_DELAY, report"`
}

func (c Config) Source() string {
	if c.SourceID != "" {
		return c.SourceID
	}
	return c.VCapApp.ApplicationID
}

type VCapApp struct {
	ApplicationID string `json:"application_id, report"`
}

func (v *VCapApp) UnmarshalEnv(jsonData string) error {
	if jsonData == "" {
		return nil
	}
	return json.Unmarshal([]byte(jsonData), &v)
}

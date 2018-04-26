package main

import "encoding/json"

type Config struct {
	LogCacheAddr string `env:"LOG_CACHE_ADDR,   required"`

	Port    int     `env:"PORT,             required"`
	VCapApp VCapApp `env:"VCAP_APPLICATION, required"`

	UAAAddr         string `env:"UAA_ADDR,          required"`
	UAAClient       string `env:"UAA_CLIENT,        required"`
	UAAClientSecret string `env:"UAA_CLIENT_SECRET, required, noreport"`

	SkipSSLValidation bool `env:"SKIP_SSL_VALIDATION"`
}

type VCapApp struct {
	ApplicationID string `json:"application_id"`
}

func (v *VCapApp) UnmarshalEnv(jsonData string) error {
	return json.Unmarshal([]byte(jsonData), &v)
}

package app

import "encoding/json"

type Drain struct {
	Namespace string `json:"namespace"`
	URL       string `json:"url"`
	All       bool   `json:"all"`
}

type Drains []Drain

func (d *Drains) UnmarshalEnv(data string) error {
	return json.Unmarshal([]byte(data), d)
}

type Config struct {
	LogsProviderAddr string `env:"LOGS_PROVIDER_ADDR, required, report"`
	LogsProviderTLS  LogsProviderTLS

	MetricsAddr string `env:"METRICS_ADDR, report"`

	ShardID string `env:"SHARD_ID, report"`
	Drains  Drains `env:"DRAINS,   report"`
}

type LogsProviderTLS struct {
	CA   string `env:"LOGS_PROVIDER_CA_FILE_PATH,   required, report"`
	Cert string `env:"LOGS_PROVIDER_CERT_FILE_PATH, required, report"`
	Key  string `env:"LOGS_PROVIDER_KEY_FILE_PATH,  required, report"`
}

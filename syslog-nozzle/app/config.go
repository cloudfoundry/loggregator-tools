package app

type Config struct {
	LogsProviderAddr string `env:"LOGS_PROVIDER_ADDR, required, report"`
	LogsProviderTLS  LogsProviderTLS

	MetricsAddr string `env:"METRICS_ADDR, report"`

	Destination string `env:"SYSLOG_DESTINATION, required, report"`
	ShardID     string `env:"SHARD_ID,                     report"`
}

type LogsProviderTLS struct {
	CA   string `env:"LOGS_PROVIDER_CA_FILE_PATH,   required, report"`
	Cert string `env:"LOGS_PROVIDER_CERT_FILE_PATH, required, report"`
	Key  string `env:"LOGS_PROVIDER_KEY_FILE_PATH,  required, report"`
}

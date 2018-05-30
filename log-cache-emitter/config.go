package main

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"io/ioutil"
	"log"

	envstruct "code.cloudfoundry.org/go-envstruct"
	"google.golang.org/grpc/credentials"
)

// Config is the configuration for a LogCache.
type Config struct {
	LogCacheAddr string `env:"LOG_CACHE_ADDR, required"`
	Addr         string `env:"ADDR,           required"`
	LogCacheTLS  TLS
}

type TLS struct {
	CAPath   string `env:"CA_PATH,   required"`
	CertPath string `env:"CERT_PATH, required"`
	KeyPath  string `env:"KEY_PATH,  required"`
}

func (t TLS) Credentials(cn string) credentials.TransportCredentials {
	creds, err := NewTLSCredentials(t.CAPath, t.CertPath, t.KeyPath, cn)
	if err != nil {
		log.Fatalf("failed to load TLS config: %s", err)
	}

	return creds
}

func NewTLSCredentials(
	caPath string,
	certPath string,
	keyPath string,
	cn string,
) (credentials.TransportCredentials, error) {
	cfg, err := NewTLSConfig(caPath, certPath, keyPath, cn)
	if err != nil {
		return nil, err
	}

	return credentials.NewTLS(cfg), nil
}

func NewTLSConfig(caPath, certPath, keyPath, cn string) (*tls.Config, error) {
	cert, err := tls.LoadX509KeyPair(certPath, keyPath)
	if err != nil {
		return nil, err
	}

	tlsConfig := &tls.Config{
		ServerName:         cn,
		Certificates:       []tls.Certificate{cert},
		InsecureSkipVerify: false,
	}

	caCertBytes, err := ioutil.ReadFile(caPath)
	if err != nil {
		return nil, err
	}

	caCertPool := x509.NewCertPool()
	if ok := caCertPool.AppendCertsFromPEM(caCertBytes); !ok {
		return nil, errors.New("cannot parse ca cert")
	}

	tlsConfig.RootCAs = caCertPool

	return tlsConfig, nil
}

// LoadConfig creates Config object from environment variables
func LoadConfig() (*Config, error) {
	c := Config{
		LogCacheAddr: ":8080",
	}

	if err := envstruct.Load(&c); err != nil {
		return nil, err
	}

	return &c, nil
}

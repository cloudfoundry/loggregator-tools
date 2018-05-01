package main

import (
	"context"
	"crypto/tls"
	"expvar"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	envstruct "code.cloudfoundry.org/go-envstruct"
	logcache "code.cloudfoundry.org/go-log-cache"
	"code.cloudfoundry.org/loggregator-tools/log-cache-forwarders/cmd/space_syslog/internal/logcacheutil"
	"code.cloudfoundry.org/loggregator-tools/log-cache-forwarders/pkg/egress/syslog"
	"code.cloudfoundry.org/loggregator-tools/log-cache-forwarders/pkg/expvarfilter"
	"code.cloudfoundry.org/loggregator-tools/log-cache-forwarders/pkg/groupmanager"
	"code.cloudfoundry.org/loggregator-tools/log-cache-forwarders/pkg/metrics"
	"code.cloudfoundry.org/loggregator-tools/log-cache-forwarders/pkg/sourceidprovider"
)

const metricsNamespace = "SpaceSyslogForwarder"

func main() {
	rand.Seed(time.Now().UnixNano())

	cfg := LoadConfig()
	envstruct.WriteReport(&cfg)

	m := startMetricsEmit()

	groupClient := logcache.NewShardGroupReaderClient(
		cfg.LogCacheHTTPAddr,
		logcache.WithHTTPClient(newOauth2HTTPClient(cfg)),
	)

	groupProvider := createGroupProvider(cfg)

	groupmanager.Start(
		cfg.LogCacheGroupName,
		time.Tick(30*time.Second),
		groupProvider,
		groupClient,
	)

	envTypes, err := logcacheutil.DrainTypeToEnvelopeTypes(cfg.DrainType)
	if err != nil {
		log.Fatalf("cannot start space syslog forwarder: %s", err)
	}

	logcache.Walk(
		context.Background(),
		cfg.LogCacheGroupName,
		syslog.NewVisitor(createSyslogWriter(cfg), m),
		groupClient.BuildReader(rand.Uint64()),
		logcache.WithWalkStartTime(time.Now()),
		logcache.WithWalkBackoff(logcache.NewAlwaysRetryBackoff(250*time.Millisecond)),
		logcache.WithWalkLimit(1000),
		logcache.WithWalkLogger(log.New(os.Stderr, "", log.LstdFlags)),
		logcache.WithWalkEnvelopeTypes(envTypes...),
	)
}

func startMetricsEmit() *metrics.Metrics {
	m := metrics.New(expvar.NewMap(metricsNamespace))

	mh := expvarfilter.NewHandler(expvar.Handler(), []string{metricsNamespace})
	go func() {
		// health endpoints (expvar)
		log.Printf("Health: %s", http.ListenAndServe(":"+os.Getenv("PORT"), mh))
	}()

	return m
}

func createGroupProvider(cfg Config) *sourceidprovider.SpaceProvider {
	return sourceidprovider.Space(
		cAPICurler{
			client: newOauth2HTTPClient(cfg),
		},
		cfg.CAPIURL,
		cfg.SpaceGUID,
	)
}

func createSyslogWriter(cfg Config) syslog.WriteCloser {
	netConf := syslog.NetworkConfig{
		Keepalive:      cfg.KeepAlive,
		DialTimeout:    cfg.DialTimeout,
		WriteTimeout:   cfg.IOTimeout,
		SkipCertVerify: cfg.SkipCertVerify,
	}
	return syslog.NewWriter(cfg.SourceHostname, cfg.SyslogURL, netConf)
}

func newOauth2HTTPClient(cfg Config) *logcache.Oauth2HTTPClient {
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: cfg.SkipCertVerify,
			},
		},
		Timeout: 5 * time.Second,
	}

	return logcache.NewOauth2HTTPClient(
		cfg.UAAURL,
		cfg.ClientID,
		cfg.ClientSecret,
		logcache.WithOauth2HTTPClient(client),
		logcache.WithOauth2HTTPUser(cfg.Username, cfg.Password),
	)
}

type cAPICurler struct {
	client *logcache.Oauth2HTTPClient
}

func (c cAPICurler) Get(url string) []byte {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Printf("Failed to create request: %s", err)
		return nil
	}

	resp, err := c.client.Do(req)
	if err != nil {
		log.Printf("Failed to make request to CAPI: %s", err)
		return nil
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Failed to read response: %s", err)
		return nil
	}

	return respBody
}

package main

import (
	"context"
	"crypto/tls"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"

	envstruct "code.cloudfoundry.org/go-envstruct"
	logcache "code.cloudfoundry.org/go-log-cache"
	"code.cloudfoundry.org/loggregator-tools/log-cache-forwarders/pkg/egress/datadog"
	"code.cloudfoundry.org/loggregator-tools/log-cache-forwarders/pkg/groupmanager"
	"code.cloudfoundry.org/loggregator-tools/log-cache-forwarders/pkg/sourceidprovider"
	datadogapi "github.com/zorkian/go-datadog-api"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	cfg := LoadConfig()
	envstruct.WriteReport(&cfg)

	groupClient := logcache.NewShardGroupReaderClient(
		cfg.LogCacheHTTPAddr,
		logcache.WithHTTPClient(newOauth2HTTPClient(cfg)),
	)

	client := logcache.NewClient(
		cfg.LogCacheHTTPAddr,
		logcache.WithHTTPClient(newOauth2HTTPClient(cfg)),
	)

	if cfg.SourceIDWhitelist != "" && cfg.SourceIDBlacklist != "" {
		log.Fatalf("Can't have a whitelist and a blacklist...")
	}

	var provider groupmanager.GroupProvider = sourceidprovider.All(client)

	if cfg.SourceIDBlacklist != "" {
		provider = sourceidprovider.NewRegex(
			true,
			cfg.SourceIDBlacklist,
			client,
		)
	}

	if cfg.SourceIDWhitelist != "" {
		provider = sourceidprovider.NewRegex(
			false,
			cfg.SourceIDWhitelist,
			client,
		)
	}

	groupmanager.Start(
		cfg.LogCacheGroupName,
		time.Tick(30*time.Second),
		provider,
		groupClient,
	)

	ddc := datadogapi.NewClient(cfg.DatadogAPIKey, "")
	visitor := datadog.Visitor(ddc, cfg.MetricHost, strings.Split(cfg.DatadogTags, ","))

	reader := groupClient.BuildReader(rand.Uint64())

	logcache.Walk(
		context.Background(),
		cfg.LogCacheGroupName,
		logcache.Visitor(visitor),
		reader,
		logcache.WithWalkStartTime(time.Now()),
		logcache.WithWalkBackoff(logcache.NewAlwaysRetryBackoff(250*time.Millisecond)),
		logcache.WithWalkLogger(log.New(os.Stderr, "", log.LstdFlags)),
	)
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
		cfg.UAAAddr,
		cfg.ClientID,
		cfg.ClientSecret,
		logcache.WithOauth2HTTPClient(client),
	)
}

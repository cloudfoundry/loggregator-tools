package main

import (
	"log"
	"os"

	envstruct "code.cloudfoundry.org/go-envstruct"
	loggregator "code.cloudfoundry.org/go-loggregator"
	"code.cloudfoundry.org/go-loggregator/rpc/loggregator_v2"
	orchestrator "code.cloudfoundry.org/go-orchestrator"
	"code.cloudfoundry.org/loggregator-tools/syslog-forwarder/internal/egress"
	"code.cloudfoundry.org/loggregator-tools/syslog-forwarder/internal/stream"
)

func main() {
	l := log.New(os.Stderr, "[Syslog-Forwarder] ", log.LstdFlags)
	l.Println("Starting Syslog Forwarder...")
	defer l.Println("Closing Syslog Forwarder...")

	cfg := LoadConfig()
	envstruct.WriteReport(&cfg)

	client := loggregator.NewRLPGatewayClient(cfg.Vcap.RLPAddr,
		loggregator.WithRLPGatewayClientLogger(l),
	)

	streamAggregator := stream.NewAggregator(client, cfg.ShardID, l)
	o := createOrchestrator(streamAggregator)

	excludeSelf := func(sourceID string) bool { return sourceID == cfg.Vcap.AppID }
	sm := stream.NewSourceManager(
		stream.NewSingleOrSpaceProvider(
			cfg.SourceID,
			cfg.Vcap.API,
			cfg.Vcap.SpaceGUID,
			cfg.IncludeServices,
			stream.WithSourceProviderSpaceExcludeFilter(excludeSelf),
		),
		o,
		cfg.UpdateInterval,
	)
	go sm.Start()

	envs := streamAggregator.Consume()
	w := createSyslogWriter(cfg, l)
	for e := range envs {
		err := w.Write(e.(*loggregator_v2.Envelope))
		if err != nil {
			l.Printf("error writing envelope to syslog: %s", err)
			continue
		}
	}
}

func createOrchestrator(s *stream.Aggregator) *orchestrator.Orchestrator {
	o := orchestrator.New(
		stream.Communicator{},
	)
	o.AddWorker(s)
	return o
}

func createSyslogWriter(cfg Config, log *log.Logger) egress.WriteCloser {
	netConf := egress.NetworkConfig{
		Keepalive:      cfg.KeepAlive,
		DialTimeout:    cfg.DialTimeout,
		WriteTimeout:   cfg.IOTimeout,
		SkipCertVerify: cfg.SkipCertVerify,
	}
	return egress.NewWriter(cfg.SourceHostname, cfg.SyslogURL, netConf, log)
}

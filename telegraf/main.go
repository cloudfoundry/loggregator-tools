package main

import (
	"fmt"
	"github.com/influxdata/toml"
	"github.com/nats-io/nats.go"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"net"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"
)

const (
	ScrapeTargetQueueName = "metrics.scrape_targets"
	appDir                = "/home/vcap/app"
	telegrafConfigDir     = appDir + "/telegraf.d"
)

var cfInstanceIP = os.Getenv("CF_INSTANCE_IP")

type TelegrafConfig struct {
	Inputs map[string]*PromInputConfig `toml:"inputs"`
}

type PromInputConfig struct {
	// An array of urls to scrape metrics from.
	URLs []string `toml:"urls"`

	TLSCA              string `toml:"tls_ca"`
	TLSCert            string `toml:"tls_cert"`
	TLSKey             string `toml:"tls_key"`
	InsecureSkipVerify bool   `toml:"insecure_skip_verify"`
}

type target struct {
	Targets []string          `json:"targets",yaml:"targets"`
	Labels  map[string]string `json:"labels",yaml:"labels"`
	Source  string            `json:"-",yaml:"source"`
}

type timestampedTarget struct {
	scrapeTarget *target
	ts           time.Time
}

type configGenerator struct {
	timestampedTargets map[string]timestampedTarget
	logger             *log.Logger
	configTTL          time.Duration
	sync.Mutex
}

func main() {
	logger := log.New(os.Stderr, "nats: ", 0)

	err := os.Mkdir(telegrafConfigDir, os.ModePerm)
	if err != nil {
		logger.Fatalf("unable to make dir(%s): %s", telegrafConfigDir, err)
	}

	natsPassword := os.Getenv("NATS_PASSWORD")
	natsServer := fmt.Sprintf("nats://nats:%s@10.0.1.0:4222", natsPassword)

	opts := nats.Options{
		Servers:           []string{natsServer},
		PingInterval:      20 * time.Second,
		AllowReconnect:    true,
		MaxReconnect:      -1,
		ReconnectWait:     100 * time.Millisecond,
		ClosedCB:          closedCB(logger),
		DisconnectedErrCB: disconnectErrHandler(logger),
		ReconnectedCB:     reconnectedCB(logger),
	}

	natsConn, err := opts.Connect()
	if err != nil {
		logger.Fatalf("Unable to connect to nats servers: %s", err)
	}

	cg := configGenerator{
		timestampedTargets: map[string]timestampedTarget{},
		logger:             logger,
		configTTL:          45 * time.Second,
	}

	_, err = natsConn.Subscribe(ScrapeTargetQueueName, cg.generate)
	if err != nil {
		logger.Fatalf("failed to subscribe to %s: %s", ScrapeTargetQueueName, err)
	}

	expirationTicker := time.NewTicker(15 * time.Second)
	writeTicker := time.NewTicker(15 * time.Second)
	for {
		select {
		case <-writeTicker.C:
			cg.writeConfigToFile()
		case <-expirationTicker.C:
			cg.expireScrapeConfigs()
		}
	}
}

func (cg *configGenerator) writeConfigToFile() {
	var urls []string
	for _, scrapeTarget := range cg.timestampedTargets {
		for _, target := range scrapeTarget.scrapeTarget.Targets {
			host, _, _ := net.SplitHostPort(target)
			if host == cfInstanceIP {
				continue
			}

			id, ok := scrapeTarget.scrapeTarget.Labels["__param_id"]
			if ok {
				urls = append(urls, fmt.Sprintf("https://%s?id=%s", target, id))
				continue
			}

			urls = append(urls, fmt.Sprintf("https://%s", target))
		}
	}

	sort.Strings(urls)

	cfg := PromInputConfig{
		URLs:               urls,
		TLSCA:              appDir + "/certs/scrape_ca.crt",
		TLSCert:            appDir + "/certs/scrape.crt",
		TLSKey:             appDir + "/certs/scrape.key",
		InsecureSkipVerify: true,
	}

	newCfgBytes, err := toml.Marshal(&TelegrafConfig{
		Inputs: map[string]*PromInputConfig{"prometheus": &cfg},
	})
	if err != nil {
		fmt.Println(err)
		return
	}

	oldCfgBytes, err := ioutil.ReadFile(telegrafConfigDir + "/inputs.conf")
	if err != nil {
		oldCfgBytes = []byte{}
	}

	if string(newCfgBytes) == string(oldCfgBytes) {
		return
	}

	pid, ok := cg.getTelegrafPID()
	if !ok {
		return
	}

	err = ioutil.WriteFile(telegrafConfigDir+"/inputs.conf", newCfgBytes, os.ModePerm)
	if err != nil {
		cg.logger.Println(err)
		return
	}

	err = syscall.Kill(pid, syscall.SIGHUP)
	if err != nil {
		cg.logger.Println(err)
	}
}

func (cg *configGenerator) getTelegrafPID() (int, bool) {
	pidBytes, err := ioutil.ReadFile(appDir + "/telegraf.pid")
	if err != nil {
		cg.logger.Println(err)
		return 0, false
	}

	pid, err := strconv.Atoi(strings.TrimSuffix(string(pidBytes), "\n"))
	if err != nil {
		cg.logger.Println(err)
		return 0, false
	}

	return pid, true
}

func (cg *configGenerator) generate(message *nats.Msg) {
	scrapeTarget, ok := cg.unmarshalScrapeTarget(message)
	if !ok {
		return
	}

	cg.addTarget(scrapeTarget)
}

func (cg *configGenerator) unmarshalScrapeTarget(message *nats.Msg) (*target, bool) {
	var t target
	err := yaml.Unmarshal(message.Data, &t)
	if err != nil {
		cg.logger.Printf("failed to unmarshal message data: %s\n", err)
		return nil, false
	}

	return &t, true
}

func (cg *configGenerator) addTarget(scrapeTarget *target) {
	cg.Lock()
	defer cg.Unlock()

	cg.timestampedTargets[scrapeTarget.Source] = timestampedTarget{
		scrapeTarget: scrapeTarget,
		ts:           time.Now(),
	}
}

func (cg *configGenerator) expireScrapeConfigs() {
	cg.Lock()
	defer cg.Unlock()

	for k, scrapeConfig := range cg.timestampedTargets {
		if time.Since(scrapeConfig.ts) >= cg.configTTL {
			delete(cg.timestampedTargets, k)
		}
	}
}

func closedCB(log *log.Logger) func(conn *nats.Conn) {
	return func(conn *nats.Conn) {
		log.Println("Nats Connection Closed")
	}
}

func reconnectedCB(log *log.Logger) func(conn *nats.Conn) {
	return func(conn *nats.Conn) {
		log.Printf("Reconnected to %s\n", conn.ConnectedUrl())
	}
}

func disconnectErrHandler(log *log.Logger) func(conn *nats.Conn, err error) {
	return func(conn *nats.Conn, err error) {
		log.Printf("Nats Error %s\n", err)
	}
}

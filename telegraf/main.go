package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"
	"github.com/influxdata/toml"
)

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

func main() {
	for {
		resp, err := http.Get("http://10.0.1.13:12345")
		if err != nil {
			fmt.Println(err)
			continue
		}

		if resp.StatusCode != http.StatusOK {
			fmt.Println(resp.StatusCode)
			continue
		}

		var urls []string
		err = json.NewDecoder(resp.Body).Decode(&urls)
		if err != nil {
			fmt.Println(err)
			continue
		}

		cfg := PromInputConfig{
			URLs:               urls,
			TLSCA:              "/home/vcap/app/certs/scrape_ca.crt",
			TLSCert:            "/home/vcap/app/certs/scrape.crt",
			TLSKey:             "/home/vcap/app/certs/scrape.key",
			InsecureSkipVerify: true,
		}

		newCfgBytes, err := toml.Marshal(&TelegrafConfig{
			Inputs: map[string]*PromInputConfig{"prometheus": &cfg},
		})
		if err != nil {
			fmt.Println(err)
			continue
		}

		err = ioutil.WriteFile("/home/vcap/app/inputs.conf", newCfgBytes, os.ModePerm)
		if err != nil {
			fmt.Println(err)
			continue
		}

		time.Sleep(time.Minute)
	}
}

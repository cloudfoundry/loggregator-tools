package processor

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"code.cloudfoundry.org/rfc5424"
)

const (
	datadogAPIEndpoint = "/api/v1/series"
)

var (
	template = `{
		"series": [{
			"metric": "%s.%s",
			"points": [[%d, %s]],
			"type": "gauge",
			"host": "%s"
		}]
	}`
)

// Getter is used to get rfc5424 encoded bytes to send to datadog.
type Getter func() ([]byte, bool)

// HTTPClient is used to communicate to the datadog API.
type HTTPClient interface {
	Post(url, contentType string, body io.Reader) (*http.Response, error)
}

// Processor gets rfc5424 messages and sends them to the datadog API.
type Processor struct {
	getter Getter
	client HTTPClient
	apiURL *url.URL
}

// New creates a new Processor.
func New(g Getter, c HTTPClient, apiBaseURL, apiKey string) *Processor {
	apiURL, err := url.Parse(apiBaseURL + datadogAPIEndpoint)
	if err != nil {
		log.Fatalf("Failed to parse datadog URL: %s", err)
	}
	query := url.Values{
		"api_key": []string{apiKey},
	}
	apiURL.RawQuery = query.Encode()

	return &Processor{
		getter: g,
		client: c,
		apiURL: apiURL,
	}
}

// Run reads from the Getter and writes to datadog. It blocks while reading.
func (p *Processor) Run() {
	for {
		data, ok := p.getter()
		if !ok {
			time.Sleep(10 * time.Millisecond)
			continue
		}

		var msg rfc5424.Message
		err := msg.UnmarshalBinary(data)
		if err != nil {
			continue
		}

		for _, sd := range msg.StructuredData {
			if strings.HasPrefix(sd.ID, "gauge@") {
				err := p.postGauge(msg, sd)
				if err != nil {
					log.Printf("failed to post to datadog: %s", err)
				}
				break
			}
		}
	}
}

func (p *Processor) postGauge(msg rfc5424.Message, sd rfc5424.StructuredData) error {
	var name, value string
	for _, p := range sd.Parameters {
		switch p.Name {
		case "name":
			name = p.Value
		case "value":
			value = p.Value
		}
	}
	body := strings.NewReader(
		fmt.Sprintf(
			template,
			msg.Hostname,
			name,
			msg.Timestamp.Unix(),
			value,
			msg.Hostname,
		),
	)
	resp, err := p.client.Post(p.apiURL.String(), "application/json", body)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode > 299 || resp.StatusCode < 200 {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("expected success status code, got %d", resp.StatusCode)
		}

		return fmt.Errorf("expected success status code, got %d: %s", resp.StatusCode, body)
	}

	return nil
}

package processor_test

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
	"sync"
	"time"

	"code.cloudfoundry.org/loggregator-tools/syslog_to_datadog/internal/processor"
	"code.cloudfoundry.org/rfc5424"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Processor", func() {
	It("posts metrics as rfc5424 message to datadog", func() {
		var getterCalled bool
		getter := func() ([]byte, bool) {
			if getterCalled {
				return nil, false
			}
			getterCalled = true
			return buildGaugeMessage(), true
		}
		spyClient := &spyClient{}
		p := processor.New(getter, spyClient, "https://an-api-addr", "an-api-key")

		go p.Run()

		Eventually(spyClient.posts).Should(HaveLen(1))
		post := spyClient.posts()[0]
		Expect(post.url).To(Equal("https://an-api-addr/api/v1/series?api_key=an-api-key"))
		Expect(post.contentType).To(Equal("application/json"))
		Expect(post.body).To(MatchJSON(`{
			"series": [
				{
					"metric": "myhostname.cpu",
					"points": [[0, 0.23]],
					"type": "gauge",
					"host": "myhostname",
					"tags": ["instance_id:4"]
				}
			]
		}`))
	})

	It("ignores messages that are not proper rfc5424 messages", func() {
		var getterCalled bool
		getter := func() ([]byte, bool) {
			if getterCalled {
				return nil, false
			}
			getterCalled = true
			return []byte("bunch o garbage"), true
		}
		spyClient := &spyClient{}
		p := processor.New(getter, spyClient, "", "an-api-key")

		go p.Run()

		Consistently(spyClient.posts).Should(BeEmpty())
	})

	It("processes counter metrics", func() {
		var getterCalled bool
		getter := func() ([]byte, bool) {
			if getterCalled {
				return nil, false
			}
			getterCalled = true
			return buildCounterMessage(), true
		}
		spyClient := &spyClient{}
		p := processor.New(getter, spyClient, "https://an-api-addr", "an-api-key")

		go p.Run()

		Eventually(spyClient.posts).Should(HaveLen(1))
		post := spyClient.posts()[0]
		Expect(post.url).To(Equal("https://an-api-addr/api/v1/series?api_key=an-api-key"))
		Expect(post.contentType).To(Equal("application/json"))
		Expect(post.body).To(MatchJSON(`{
			"series": [
				{
					"metric": "myhostname.requests",
					"points": [[0, 1234]],
					"type": "gauge",
					"host": "myhostname",
					"tags": ["instance_id:4"]
				}
			]
		}`))
	})

	It("only sends the first structured data element", func() {
		var getterCalled bool
		getter := func() ([]byte, bool) {
			if getterCalled {
				return nil, false
			}
			getterCalled = true
			return buildBloatedGaugeMessage(), true
		}
		spyClient := &spyClient{}
		p := processor.New(getter, spyClient, "", "an-api-key")

		go p.Run()

		Eventually(spyClient.posts).Should(HaveLen(1))
		Consistently(spyClient.posts).Should(HaveLen(1))
	})
})

func buildGaugeMessage() []byte {
	m := rfc5424.Message{
		Priority:  rfc5424.Daemon | rfc5424.Info,
		Timestamp: time.Unix(0, 0),
		Hostname:  "myhostname",
		AppName:   "someapp",
		ProcessID: "[4]",
		StructuredData: []rfc5424.StructuredData{
			{
				ID: "gauge@47450",
				Parameters: []rfc5424.SDParam{
					{
						Name:  "name",
						Value: "cpu",
					},
					{
						Name:  "value",
						Value: "0.23",
					},
					{
						Name:  "unit",
						Value: "percentage",
					},
				},
			},
		},
	}

	data, err := m.MarshalBinary()
	Expect(err).ToNot(HaveOccurred())
	return data
}

func buildCounterMessage() []byte {
	m := rfc5424.Message{
		Priority:  rfc5424.Daemon | rfc5424.Info,
		Timestamp: time.Unix(0, 0),
		Hostname:  "myhostname",
		AppName:   "someapp",
		Message:   []byte("some log"),
		ProcessID: "[4]",
		StructuredData: []rfc5424.StructuredData{
			{
				ID: "counter@47450",
				Parameters: []rfc5424.SDParam{
					{
						Name:  "name",
						Value: "requests",
					},
					{
						Name:  "total",
						Value: "1234",
					},
					{
						Name:  "delta",
						Value: "5",
					},
				},
			},
		},
	}

	data, err := m.MarshalBinary()
	Expect(err).ToNot(HaveOccurred())
	return data
}

func buildBloatedGaugeMessage() []byte {
	m := rfc5424.Message{
		Priority:  rfc5424.Daemon | rfc5424.Info,
		Timestamp: time.Unix(0, 0),
		Hostname:  "myhostname",
		AppName:   "someapp",
		StructuredData: []rfc5424.StructuredData{
			{
				ID: "gauge@47450",
				Parameters: []rfc5424.SDParam{
					{
						Name:  "name",
						Value: "cpu",
					},
					{
						Name:  "value",
						Value: "0.23",
					},
					{
						Name:  "unit",
						Value: "percentage",
					},
				},
			},
			{
				ID: "gauge@47450",
				Parameters: []rfc5424.SDParam{
					{
						Name:  "name",
						Value: "memory",
					},
					{
						Name:  "value",
						Value: "0.23",
					},
					{
						Name:  "unit",
						Value: "percentage",
					},
				},
			},
		},
	}

	data, err := m.MarshalBinary()
	Expect(err).ToNot(HaveOccurred())
	return data
}

type post struct {
	url         string
	contentType string
	body        []byte
}

type spyClient struct {
	mu     sync.Mutex
	_posts []post
}

func (s *spyClient) posts() []post {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s._posts
}

func (s *spyClient) Post(url, contentType string, body io.Reader) (*http.Response, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	buf := bytes.NewBuffer(nil)
	_, err := buf.ReadFrom(body)
	Expect(err).ToNot(HaveOccurred())

	s._posts = append(s._posts, post{
		url,
		contentType,
		buf.Bytes(),
	})

	return &http.Response{
			StatusCode: 201,
			Body:       ioutil.NopCloser(bytes.NewBuffer(nil)),
		},
		nil
}

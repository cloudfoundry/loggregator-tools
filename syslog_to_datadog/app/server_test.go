package app_test

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"

	"code.cloudfoundry.org/loggregator-tools/syslog_to_datadog/app"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Server", func() {
	It("forwards rfc5424 messages containing metrics to datadog", func() {
		datadog := newSpyDatadog()
		s := app.NewServer(":0", "junk-key", app.WithDatadogBaseURL(datadog.server.URL))
		go s.Run()

		resp, err := http.Post(fmt.Sprintf(
			"http://%s", s.Addr()),
			"text/plain",
			strings.NewReader(`<30>1 2017-10-04T13:00:52.662629-06:00 myhostname someapp [4] - [gauge@47450 name="cpu" value="0.23" unit="percentage"]`),
		)
		Expect(err).ToNot(HaveOccurred())
		Expect(resp.StatusCode).To(Equal(http.StatusAccepted))
		Eventually(datadog.requests).Should(HaveLen(1))

		req := datadog.requests()[0]
		Expect(req.body).To(MatchJSON(`{
			"series": [
				{
					"metric": "myhostname.cpu",
					"points": [[1507143652, 0.23]],
					"type": "gauge",
					"host": "myhostname",
					"tags": ["instance_id:4"]
				}
			]
		}`))
	})
})

type request struct {
	url         string
	body        string
	contentType string
}

type spyDatadog struct {
	mu        sync.Mutex
	_requests []request

	server *httptest.Server
}

func newSpyDatadog() *spyDatadog {
	s := &spyDatadog{}
	s.server = httptest.NewServer(s)

	return s
}

func (s *spyDatadog) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	buf := bytes.NewBuffer(make([]byte, 0, r.ContentLength))
	_, err := buf.ReadFrom(r.Body)
	Expect(err).ToNot(HaveOccurred())

	s.mu.Lock()
	s._requests = append(s._requests, request{
		url:         r.URL.String(),
		body:        buf.String(),
		contentType: r.Header.Get("Content-Type"),
	})
	s.mu.Unlock()

	w.WriteHeader(http.StatusAccepted)
}

func (s *spyDatadog) requests() []request {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s._requests
}

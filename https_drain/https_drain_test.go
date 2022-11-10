package main_test

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os/exec"
	"strings"
	"time"

	"code.cloudfoundry.org/rfc5424"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
)

var _ = Describe("HttpsDrain", func() {
	var (
		port = 1354
	)

	It("writes out the body of the request it receives", func() {
		_, kill, stdout, _ := startHTTPSDrain(port, "")
		defer kill()

		body := "to check body output"
		_ = writeLog(port, body)

		Expect(stdout.Contents()).To(ContainSubstring(body))
	})

	It("returns a 400 when no body is given", func() {
		_, kill, _, _ := startHTTPSDrain(port, "")
		defer kill()

		resp := writeLog(port, "")

		Expect(resp.StatusCode).To(Equal(http.StatusBadRequest))
	})

	Context("with a COUNTER_URL", func() {
		It("reports log counts to the COUNTER_URL", func() {
			reqs := make(chan *http.Request, 100)
			bodies := make(chan []byte, 100)
			ts := httptest.NewServer(
				http.HandlerFunc(
					func(w http.ResponseWriter, r *http.Request) {
						reqs <- r

						body, err := io.ReadAll(r.Body)
						if err != nil {
							panic(err)
						}

						bodies <- body
					},
				),
			)
			defer ts.Close()

			// 1s before it posts
			_, kill, _, _ := startHTTPSDrain(port, ts.URL)
			defer kill()

			var req *http.Request
			Eventually(reqs).Should(Receive(&req))
			Expect(req.Method).To(Equal("POST"))
			Expect(req.URL.RequestURI()).To(Equal("/set/"))

			rfcLog := rfc5424.Message{
				Priority:  rfc5424.Emergency,
				Timestamp: time.Now(),
				Hostname:  "some-host",
				AppName:   "some-app",
				ProcessID: "procID",
				MessageID: "msgID",
				Message:   []byte(`{"id":"some-id", "msgCount":1}`),
			}

			data, err := rfcLog.MarshalBinary()
			Expect(err).ToNot(HaveOccurred())

			writeLog(port, string(data))
			Eventually(bodies).Should(Receive(MatchJSON(`[{"id":"some-id", "msgCount":1, "primeCount":0}]`)))
		})
	})
})

func startHTTPSDrain(port int, counterUrl string) (*gexec.Session, func(), *gbytes.Buffer, *gbytes.Buffer) {
	path, err := gexec.Build("code.cloudfoundry.org/loggregator-tools/https_drain")
	Expect(err).ToNot(HaveOccurred())

	cmd := exec.Command(path)
	cmd.Env = []string{
		fmt.Sprintf("PORT=%d", port),
		fmt.Sprintf("COUNTER_URL=%s", counterUrl),
		"INTERVAL=100ms",
	}

	stdout, stderr := gbytes.NewBuffer(), gbytes.NewBuffer()
	session, err := gexec.Start(cmd, stdout, stderr)
	Expect(err).ToNot(HaveOccurred())

	return session, func() { session.Kill().Wait() }, stdout, stderr
}

func writeLog(port int, body string) *http.Response {
	// log := `107 <7>1 2016-02-28T09:57:10.804642398-05:00 myhostname someapp - - [foo@1234 Revision="1.2.3.4"] Hello, World!`
	var response *http.Response
	Eventually(func() error {
		var err error
		response, err = http.Post(
			fmt.Sprintf("http://localhost:%d", port),
			"text/plain",
			strings.NewReader(body),
		)
		return err
	}).ShouldNot(HaveOccurred())

	return response
}

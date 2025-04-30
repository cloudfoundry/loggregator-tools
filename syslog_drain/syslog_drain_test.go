package main_test

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"os/exec"
	"time"

	"code.cloudfoundry.org/rfc5424"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
)

var _ = Describe("SyslogDrain", func() {
	var (
		stdout  *gbytes.Buffer
		session *gexec.Session
		port    int
		writer  io.WriteCloser
		ts      *httptest.Server
		reqs    chan *http.Request
		bodies  chan []byte
	)

	BeforeEach(func() {
		port = 9999
		path, err := gexec.Build("code.cloudfoundry.org/loggregator-tools/syslog_drain")
		Expect(err).ToNot(HaveOccurred())

		reqs = make(chan *http.Request, 100)
		bodies = make(chan []byte, 100)
		ts = httptest.NewServer(
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
		cmd := exec.Command(path)
		cmd.Env = []string{
			fmt.Sprintf("PORT=%d", port),
			fmt.Sprintf("COUNTER_URL=%s", ts.URL),
			"INTERVAL=100ms",
		}

		stdout = gbytes.NewBuffer()
		session, err = gexec.Start(cmd, stdout, GinkgoWriter)
		Expect(err).ToNot(HaveOccurred())

		Eventually(func() error {
			writer, err = net.Dial("tcp4", fmt.Sprintf("127.0.0.1:%d", port))
			return err
		}, 3).ShouldNot(HaveOccurred())

	})

	AfterEach(func() {
		err := writer.Close()
		Expect(err).ToNot(HaveOccurred())
		session.Kill()
		ts.Close()
	})

	It("writes out the body of the request it receives", func() {
		writeRfcLog(`{"id":"some-id", "msgCount":1}`, writer)

		Eventually(stdout.Contents).Should(ContainSubstring(`{"id":"some-id", "msgCount":1}`))
	})

	It("reports syslog counts to counter", func() {
		writeRfcLog(`{"id":"some-id", "msgCount":1}`, writer)

		var req *http.Request
		Eventually(reqs).Should(Receive(&req))
		Expect(req.Method).To(Equal("POST"))
		Expect(req.URL.RequestURI()).To(Equal("/set/"))

		Eventually(bodies).Should(Receive(MatchJSON(`[{"id":"some-id", "msgCount":1, "primeCount":0}]`)))
	})
})

func writeRfcLog(msg string, writer io.Writer) {
	rfcLog := rfc5424.Message{
		Priority:  rfc5424.Emergency,
		Timestamp: time.Now(),
		Hostname:  "some-host",
		AppName:   "some-app",
		ProcessID: "procID",
		MessageID: "msgID",
		Message:   []byte(msg),
	}

	_, err := rfcLog.WriteTo(writer)
	Expect(err).ToNot(HaveOccurred())
}

package cloudcontroller_test

import (
	"errors"
	"strings"

	"code.cloudfoundry.org/cli/plugin"
	"code.cloudfoundry.org/loggregator-tools/syslog-forwarder/internal/cloudcontroller"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

// Assert that the CurlClient is a Curler
var _ cloudcontroller.Curler = &cloudcontroller.CLICurlClient{}

var _ = Describe("CurlClient", func() {
	var (
		conn *stubCliConnection
		c    *cloudcontroller.CLICurlClient
	)

	BeforeEach(func() {
		conn = newStubCliConnection()
		c = cloudcontroller.NewCLICurlClient(conn)
	})

	It("uses 'curl' command", func() {
		_, err := c.Curl("some-url", "GET", "")
		Expect(err).ToNot(HaveOccurred())
		Expect(conn.args).To(HaveLen(1))
		Expect(conn.args[0][0]).To(Equal("curl"))
		Expect(conn.args[0][1]).To(Equal("some-url"))
	})

	It("returns the response as a joined byte slice", func() {
		conn.resp["curl some-url"] = `{
			"snacks" : []
		}`
		resp, err := c.Curl("some-url", "GET", "")

		Expect(string(resp)).To(Equal(`{
			"snacks" : []
		}`))
		Expect(err).ToNot(HaveOccurred())
	})

	It("returns any error", func() {
		conn.err = errors.New("some error")
		_, err := c.Curl("some-url", "GET", "")

		Expect(err).To(HaveOccurred())
	})

	It("panics if method is not GET or body is populated", func() {
		Expect(func() {
			_, _ = c.Curl("some-url", "POST", "")
		}).To(Panic())

		Expect(func() {
			_, _ = c.Curl("some-url", "GET", "something")
		}).To(Panic())
	})
})

type stubCliConnection struct {
	plugin.CliConnection

	args [][]string
	err  error
	resp map[string]string
}

func newStubCliConnection() *stubCliConnection {
	return &stubCliConnection{
		resp: make(map[string]string),
	}
}

func (s *stubCliConnection) CliCommandWithoutTerminalOutput(args ...string) ([]string, error) {
	s.args = append(
		s.args,
		args,
	)

	output := s.resp[strings.Join(args, " ")]
	return strings.Split(output, "\n"), s.err
}

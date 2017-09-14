package main_test

import (
	"bytes"
	"fmt"
	"net/http"
	"net/url"
	"os/exec"

	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Logspinner", func() {
	var port = 9423

	It("defaults write delay to 1s", func() {
		_, kill, stdout, _ := startLogSpinner(port)
		defer kill()

		buffer := requestToStartLogs(port, url.Values{
			"cycles": {"1"},
		})

		Expect(buffer.String()).To(Equal("cycles 1, delay 1s, text LogSpinner Log Message\n"))

		Eventually(stdout.Contents, 2).Should(ContainSubstring("Duration"))

		lines := bytes.Split(stdout.Contents(), []byte("\n"))
		Expect(lines).To(HaveLen(4))
		Expect(string(lines[1])).To(Equal("msg 1 LogSpinner Log Message"))
	})

	It("defaults write cycles to 10", func() {
		_, kill, stdout, _ := startLogSpinner(port)
		defer kill()

		buffer := requestToStartLogs(port, url.Values{
			"delay": {"1ms"},
		})

		Expect(buffer.String()).To(Equal("cycles 10, delay 1ms, text LogSpinner Log Message\n"))

		Eventually(stdout.Contents).Should(ContainSubstring("Duration"))

		lines := bytes.Split(stdout.Contents(), []byte("\n"))
		Expect(lines).To(HaveLen(13))
		Expect(string(lines[1])).To(Equal("msg 1 LogSpinner Log Message"))
		Expect(string(lines[10])).To(Equal("msg 10 LogSpinner Log Message"))
	})

	It("accepts with a custom log message", func() {
		_, kill, stdout, _ := startLogSpinner(port)
		defer kill()

		buffer := requestToStartLogs(port, url.Values{
			"delay": {"1ms"},
			"text":  {"Custom"},
		})

		Expect(buffer.String()).To(Equal("cycles 10, delay 1ms, text Custom\n"))

		Eventually(stdout.Contents).Should(ContainSubstring("Duration"))

		lines := bytes.Split(stdout.Contents(), []byte("\n"))
		Expect(lines).To(HaveLen(13))
		Expect(string(lines[1])).To(Equal("msg 1 Custom"))
		Expect(string(lines[10])).To(Equal("msg 10 Custom"))
	})
})

func startLogSpinner(port int) (*gexec.Session, func(), *gbytes.Buffer, *gbytes.Buffer) {
	path, err := gexec.Build("code.cloudfoundry.org/loggregator-tools/logspinner")
	Expect(err).ToNot(HaveOccurred())

	cmd := exec.Command(path)
	cmd.Env = []string{fmt.Sprintf("PORT=%d", port)}

	stdout, stderr := gbytes.NewBuffer(), gbytes.NewBuffer()
	session, err := gexec.Start(cmd, stdout, stderr)
	Expect(err).ToNot(HaveOccurred())

	return session, func() { session.Kill().Wait() }, stdout, stderr
}

func requestToStartLogs(port int, query url.Values) *bytes.Buffer {
	var resp *http.Response
	Eventually(func() error {
		var err error
		resp, err = http.Get(fmt.Sprintf("http://localhost:%d?%s", port, query.Encode()))

		return err
	}).ShouldNot(HaveOccurred())
	defer resp.Body.Close()

	buffer := bytes.NewBuffer(nil)
	buffer.ReadFrom(resp.Body)

	return buffer
}

package egress_test

import (
	"net/url"

	"code.cloudfoundry.org/loggregator-tools/log-cache-forwarders/cmd/syslog/internal/egress"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("WriterFactory", func() {
	It("returns an https writer when the url begins with http", func() {
		url, err := url.Parse("https://the-syslog-endpoint.com")
		Expect(err).ToNot(HaveOccurred())

		writer := egress.NewWriter("source-id", "source-host", url, egress.NetworkConfig{})

		_, ok := writer.(*egress.HTTPSWriter)
		Expect(ok).To(BeTrue())
	})

	It("returns a tcp writer when the url begins with syslog://", func() {
		url, err := url.Parse("syslog://the-syslog-endpoint.com")
		Expect(err).ToNot(HaveOccurred())

		writer := egress.NewWriter("source-id", "source-host", url, egress.NetworkConfig{})

		_, ok := writer.(*egress.TCPWriter)
		Expect(ok).To(BeTrue())
	})
})

package syslog_test

import (
	"net/url"

	"code.cloudfoundry.org/loggregator-tools/log-cache-forwarders/pkg/egress/syslog"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("WriterFactory", func() {
	It("returns an https writer when the url begins with http", func() {
		url, err := url.Parse("https://the-syslog-endpoint.com")
		Expect(err).ToNot(HaveOccurred())

		writer := syslog.NewWriter("source-host", url, syslog.NetworkConfig{})

		_, ok := writer.(*syslog.HTTPSWriter)
		Expect(ok).To(BeTrue())
	})

	It("returns a tcp writer when the url begins with syslog://", func() {
		url, err := url.Parse("syslog://the-syslog-endpoint.com")
		Expect(err).ToNot(HaveOccurred())

		writer := syslog.NewWriter("source-host", url, syslog.NetworkConfig{})

		_, ok := writer.(*syslog.TCPWriter)
		Expect(ok).To(BeTrue())
	})
})

package config_test

import (
	"code.cloudfoundry.org/loggregator-tools/syslog-forwarder/internal/testhelper"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"os"

	"code.cloudfoundry.org/loggregator-tools/syslog-forwarder/internal/egress/config"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("TLS", func() {

	Context("NewMutalTLSConfig", func() {

		It("builds a config struct", func() {
			conf, err := config.NewMutualTLSConfig(
				testhelper.Cert("localhost.crt"),
				testhelper.Cert("localhost.key"),
				testhelper.Cert("loggregator-ca.crt"),
				"test-server-name",
			)
			Expect(err).ToNot(HaveOccurred())

			Expect(conf.Certificates).To(HaveLen(1))
			Expect(conf.InsecureSkipVerify).To(BeFalse())
			Expect(conf.ClientAuth).To(Equal(tls.RequireAndVerifyClientCert))
			Expect(conf.MinVersion).To(Equal(uint16(tls.VersionTLS12)))
			Expect(conf.CipherSuites).To(ConsistOf(
				tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
				tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			))

			Expect(string(conf.RootCAs.Subjects()[0])).To(ContainSubstring("loggregatorCA"))
			Expect(string(conf.ClientCAs.Subjects()[0])).To(ContainSubstring("loggregatorCA"))

			Expect(conf.ServerName).To(Equal("test-server-name"))
		})

		It("allows you to not specify a CA cert", func() {
			conf, err := config.NewMutualTLSConfig(
				testhelper.Cert("localhost.crt"),
				testhelper.Cert("localhost.key"),
				"",
				"",
			)
			Expect(err).ToNot(HaveOccurred())

			Expect(conf.RootCAs).To(BeNil())
			Expect(conf.ClientCAs).To(BeNil())
		})

		It("returns an error when given invalid cert/key paths", func() {
			_, err := config.NewMutualTLSConfig("", "", testhelper.Cert("loggregator-ca.crt"), "")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("failed to load keypair: open : no such file or directory"))
		})

		It("returns an error when given invalid ca cert path", func() {
			_, err := config.NewMutualTLSConfig(testhelper.Cert("localhost.crt"), testhelper.Cert("localhost.key"), "/file/that/does/not/exist", "")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("failed to read ca cert file: open /file/that/does/not/exist: no such file or directory"))
		})

		It("returns an error when given invalid ca cert file", func() {
			empty := writeFile("")
			defer func() {
				err := os.Remove(empty)
				Expect(err).ToNot(HaveOccurred())
			}()
			_, err := config.NewMutualTLSConfig(testhelper.Cert("localhost.crt"), testhelper.Cert("localhost.key"), empty, "")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("unable to load ca cert file"))
		})

		It("returns an error when the certificate is not signed by the CA", func() {
			_, err := config.NewMutualTLSConfig(testhelper.Cert("localhost.crt"), testhelper.Cert("localhost.key"), testhelper.Cert("non-signing-ca.crt"), "")
			Expect(err).To(HaveOccurred())
			_, ok := err.(config.CASignatureError)
			Expect(ok).To(BeTrue())
		})
	})

	Context("NewTLSConfig", func() {
		It("returns basic TLS config", func() {
			tlsConf := config.NewTLSConfig()
			Expect(tlsConf.InsecureSkipVerify).To(BeFalse())
			Expect(tlsConf.ClientAuth).To(Equal(tls.NoClientCert))
			Expect(tlsConf.MinVersion).To(Equal(uint16(tls.VersionTLS12)))
			Expect(tlsConf.CipherSuites).To(ContainElement(tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256))
			Expect(tlsConf.CipherSuites).To(ContainElement(tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384))
		})
	})
})

func writeFile(data string) string {
	f, err := ioutil.TempFile("", "")
	Expect(err).ToNot(HaveOccurred())
	_, err = fmt.Fprintf(f, data)
	Expect(err).ToNot(HaveOccurred())
	return f.Name()
}
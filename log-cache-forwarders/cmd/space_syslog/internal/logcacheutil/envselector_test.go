package logcacheutil_test

import (
	"code.cloudfoundry.org/go-log-cache/rpc/logcache_v1"
	"code.cloudfoundry.org/loggregator-tools/log-cache-forwarders/cmd/space_syslog/internal/logcacheutil"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Envselector", func() {
	It("returns LOG for drain type logs", func() {
		Expect(logcacheutil.DrainTypeToEnvelopeTypes("logs")).To(ConsistOf(
			logcache_v1.EnvelopeType_LOG,
		))
	})

	It("returns COUNTER and GAUGE for drain type metrics", func() {
		Expect(logcacheutil.DrainTypeToEnvelopeTypes("metrics")).To(
			ConsistOf(
				logcache_v1.EnvelopeType_COUNTER,
				logcache_v1.EnvelopeType_GAUGE,
			))
	})

	It("returns COUNTER, GAUGE, and LOG for drain type all", func() {
		Expect(logcacheutil.DrainTypeToEnvelopeTypes("all")).To(
			ConsistOf(
				logcache_v1.EnvelopeType_LOG,
				logcache_v1.EnvelopeType_COUNTER,
				logcache_v1.EnvelopeType_GAUGE,
			))
	})

	It("returns an error when an unknown drain type is passed", func() {
		_, err := logcacheutil.DrainTypeToEnvelopeTypes("foo")
		Expect(err).To(MatchError("unknown drain type: foo"))
	})
})

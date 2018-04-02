package sourceidprovider_test

import (
	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	rpc "code.cloudfoundry.org/go-log-cache/rpc/logcache_v1"
	"code.cloudfoundry.org/loggregator-tools/log-cache-forwarders/datadog/internal/sourceidprovider"
)

var _ = Describe("AllProvider", func() {
	var stubMetaFetcher *stubMetaFetcher

	BeforeEach(func() {
		stubMetaFetcher = newStubMetaFetcher()
	})

	It("gets the MetaInfo and returns the filtered sourceIDs", func() {
		stubMetaFetcher.metaResponse = map[string]*rpc.MetaInfo{
			"88d17322-37f6-4b1f-8e3d-9e4a20b04efe": {},
			"doppler": {},
		}

		provider := sourceidprovider.All(stubMetaFetcher)

		Expect(provider.SourceIDs()).To(ConsistOf(
			"88d17322-37f6-4b1f-8e3d-9e4a20b04efe",
			"doppler",
		))
	})

	It("returns an empty slice when meta has an error", func() {
		stubMetaFetcher.metaError = errors.New("meta had an issue")

		provider := sourceidprovider.All(stubMetaFetcher)

		Expect(provider.SourceIDs()).To(HaveLen(0))
	})

	It("ignores empty source id", func() {
		stubMetaFetcher.metaResponse = map[string]*rpc.MetaInfo{
			"": {},
		}

		provider := sourceidprovider.All(stubMetaFetcher)

		Expect(provider.SourceIDs()).To(HaveLen(0))
	})
})

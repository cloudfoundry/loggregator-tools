package sourceidprovider_test

import (
	"context"
	"errors"

	rpc "code.cloudfoundry.org/go-log-cache/rpc/logcache_v1"
	"code.cloudfoundry.org/loggregator-tools/log-cache-forwarders/pkg/sourceidprovider"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

const guidRegex = `[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}`

var _ = Describe("RegexProvdier", func() {
	var (
		stubMetaFetcher *stubMetaFetcher
	)

	BeforeEach(func() {
		stubMetaFetcher = newStubMetaFetcher()
	})

	It("gets the MetaInfo and returns the filtered sourceIDs", func() {
		stubMetaFetcher.metaResponse = map[string]*rpc.MetaInfo{
			"88d17322-37f6-4b1f-8e3d-9e4a20b04efe": {},
			"doppler": {},
		}

		whitelistProvider := sourceidprovider.NewRegex(
			false,
			guidRegex,
			stubMetaFetcher,
		)
		Expect(whitelistProvider.SourceIDs()).To(ConsistOf("88d17322-37f6-4b1f-8e3d-9e4a20b04efe"))

		blacklistProvider := sourceidprovider.NewRegex(
			true,
			guidRegex,
			stubMetaFetcher,
		)
		Expect(blacklistProvider.SourceIDs()).To(ConsistOf("doppler"))
	})

	It("returns an empty slice when meta has an error", func() {
		stubMetaFetcher.metaError = errors.New("meta had an issue")

		blacklistProvider := sourceidprovider.NewRegex(
			true,
			guidRegex,
			stubMetaFetcher,
		)
		Expect(blacklistProvider.SourceIDs()).To(HaveLen(0))
	})
})

func newStubMetaFetcher() *stubMetaFetcher {
	return &stubMetaFetcher{}
}

type stubMetaFetcher struct {
	metaError    error
	metaResponse map[string]*rpc.MetaInfo
}

func (s *stubMetaFetcher) Meta(ctx context.Context) (map[string]*rpc.MetaInfo, error) {
	return s.metaResponse, s.metaError
}

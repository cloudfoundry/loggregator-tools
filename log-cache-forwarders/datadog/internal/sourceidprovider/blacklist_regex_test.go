package sourceidprovider_test

import (
	"context"
	"errors"

	rpc "code.cloudfoundry.org/go-log-cache/rpc/logcache"
	"code.cloudfoundry.org/loggregator-tools/log-cache-forwarders/datadog/internal/sourceidprovider"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("BlacklistRegexProvdier", func() {
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

		regexProvider := sourceidprovider.NewBlacklistRegex("[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}", stubMetaFetcher)

		Expect(regexProvider.SourceIDs()).To(ConsistOf("doppler"))
	})

	It("returns an empty string when meta has an error", func() {
		stubMetaFetcher.metaError = errors.New("meta had an issue")

		regexProvider := sourceidprovider.NewBlacklistRegex("[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}", stubMetaFetcher)

		Expect(regexProvider.SourceIDs()).To(HaveLen(0))
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

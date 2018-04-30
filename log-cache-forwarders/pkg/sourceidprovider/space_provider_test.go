package sourceidprovider_test

import (
	"code.cloudfoundry.org/loggregator-tools/log-cache-forwarders/pkg/sourceidprovider"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("SpaceProvider", func() {
	It("gets all apps from the space", func() {
		curler := &stubCurler{
			resp: []byte(`
				{
					"resources": [
						{"guid": "app-1"},
						{"guid": "app-2"},
						{"guid": "app-3"}
					]
				}
				`),
		}

		provider := sourceidprovider.Space(curler, "http://hostname.com", "space-guid")

		Expect(provider.SourceIDs()).To(ConsistOf(
			"app-1",
			"app-2",
			"app-3",
		))

		Expect(curler.requestedURL).To(Equal("http://hostname.com/v3/apps?space_guids=space-guid"))
	})

})

type stubCurler struct {
	resp         []byte
	requestedURL string
}

func (s *stubCurler) Get(url string) []byte {
	s.requestedURL = url
	return s.resp
}

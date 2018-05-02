package sourceidprovider_test

import (
	"code.cloudfoundry.org/loggregator-tools/log-cache-forwarders/pkg/sourceidprovider"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("SpaceProvider", func() {
	It("gets all apps from the space", func() {
		curler := &stubCurler{
			resp: [][]byte{
				[]byte(`
				{
					"resources": [
						{"guid": "app-1"},
						{"guid": "app-2"},
						{"guid": "app-3"}
					]
				}`),
			},
		}

		provider := sourceidprovider.Space(curler, "http://hostname.com", "space-guid")

		Expect(provider.SourceIDs()).To(ConsistOf(
			"app-1",
			"app-2",
			"app-3",
		))

		Expect(curler.requestedURLs).To(Equal([]string{
			"http://hostname.com/v3/apps?space_guids=space-guid",
		}))
	})

	It("can get services in addition to apps", func() {
		curler := &stubCurler{
			resp: [][]byte{
				[]byte(`
				{
					"resources": [
						{"guid": "app-1"},
						{"guid": "app-2"},
						{"guid": "app-3"}
					]
				}`),
				[]byte(`{
					"resources": [
						{"guid": "service-1"},
						{"guid": "service-2"},
						{"guid": "service-3"}
					]
				}`),
			},
		}

		provider := sourceidprovider.Space(
			curler,
			"http://hostname.com",
			"space-guid",
			sourceidprovider.WithSpaceServiceInstances(),
		)

		Expect(provider.SourceIDs()).To(ConsistOf(
			"app-1",
			"app-2",
			"app-3",
			"service-1",
			"service-2",
			"service-3",
		))

		Expect(curler.requestedURLs).To(Equal([]string{
			"http://hostname.com/v3/apps?space_guids=space-guid",
			"http://hostname.com/v3/service_instances?space_guids=space-guid",
		}))
	})
})

type stubCurler struct {
	requestCount  int
	resp          [][]byte
	requestedURLs []string
}

func (s *stubCurler) Get(url string) []byte {
	s.requestedURLs = append(s.requestedURLs, url)
	resp := s.resp[s.requestCount]
	s.requestCount++

	return resp
}

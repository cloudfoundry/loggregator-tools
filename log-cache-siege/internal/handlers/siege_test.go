package handlers_test

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"

	"code.cloudfoundry.org/go-log-cache/rpc/logcache"
	"code.cloudfoundry.org/loggregator-tools/log-cache-siege/internal/handlers"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Siege", func() {
	var (
		recorder       *httptest.ResponseRecorder
		s              http.Handler
		spyMetaFetcher *spyMetaFetcher
		server         *httptest.Server

		requests    chan string
		statusCodes chan int
	)

	BeforeEach(func() {
		requests = make(chan string, 10)
		statusCodes = make(chan int, 10)
		server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requests <- fmt.Sprintf("%s->%s", r.Method, r.URL.Path)

			select {
			case code := <-statusCodes:
				w.WriteHeader(code)
			default:
			}
		}))

		spyMetaFetcher = newSpyMetaFetcher()
		recorder = httptest.NewRecorder()
		s = handlers.NewSiege(server.URL, http.DefaultClient, spyMetaFetcher)
	})

	It("starts a request spinner for each source ID from meta", func() {
		spyMetaFetcher.results = map[string]*logcache.MetaInfo{
			"a": &logcache.MetaInfo{},
			"b": &logcache.MetaInfo{},
			"c": &logcache.MetaInfo{},
		}

		req := httptest.NewRequest("POST", "/v1/siege", nil)
		s.ServeHTTP(recorder, req)

		Expect(recorder.Code).To(Equal(http.StatusOK))

		Eventually(requests).Should(HaveLen(3))

		var rs []string
		for i := 0; i < 3; i++ {
			rs = append(rs, <-requests)
		}

		Expect(rs).To(ConsistOf(
			"POST->/v1/start/a",
			"POST->/v1/start/b",
			"POST->/v1/start/c",
		))
	})

	It("returns a 500 when the request to request spinner fails", func() {
		failer := &httpFailer{}
		spyMetaFetcher.results = map[string]*logcache.MetaInfo{
			"a": &logcache.MetaInfo{},
		}

		s = handlers.NewSiege(server.URL, failer, spyMetaFetcher)
		req := httptest.NewRequest("POST", "/v1/siege", nil)
		s.ServeHTTP(recorder, req)

		Expect(recorder.Code).To(Equal(http.StatusInternalServerError))
		Expect(recorder.Body.String()).To(MatchJSON(`{"error":"some-error"}`))
	})

	It("returns a 500 when the request spinner returns non 200", func() {
		spyMetaFetcher.results = map[string]*logcache.MetaInfo{
			"a": &logcache.MetaInfo{},
		}
		statusCodes <- 410
		req := httptest.NewRequest("POST", "/v1/siege", nil)
		s.ServeHTTP(recorder, req)

		Expect(recorder.Code).To(Equal(http.StatusInternalServerError))
		Expect(recorder.Body.String()).To(MatchJSON(`{"error":"unexpected status code from request spinner: 410"}`))
	})

	It("uses the request context to fetch the meta", func() {
		ctx, cancel := context.WithCancel(context.Background())

		req := httptest.NewRequest("POST", "/v1/siege", nil)
		req = req.WithContext(ctx)
		s.ServeHTTP(recorder, req)

		cancel()
		Expect(spyMetaFetcher.ctx.Done()).To(BeClosed())
	})

	It("returns a 500 when the meta fetcher fails", func() {
		spyMetaFetcher.err = errors.New("some-error")
		req := httptest.NewRequest("POST", "/v1/siege", nil)
		s.ServeHTTP(recorder, req)

		Expect(recorder.Code).To(Equal(http.StatusInternalServerError))
		Expect(recorder.Body.String()).To(MatchJSON(`{"error":"some-error"}`))
	})

	It("returns a 405 for a non-POST", func() {
		req := httptest.NewRequest("GET", "/v1/siege", nil)
		s.ServeHTTP(recorder, req)

		Expect(recorder.Code).To(Equal(http.StatusMethodNotAllowed))
	})

	It("returns 404 for non /v1/siege endpoint", func() {
		req := httptest.NewRequest("POST", "/v1/invalid", nil)
		s.ServeHTTP(recorder, req)

		Expect(recorder.Code).To(Equal(http.StatusNotFound))
	})
})

type spyMetaFetcher struct {
	ctx     context.Context
	results map[string]*logcache.MetaInfo
	err     error
}

func newSpyMetaFetcher() *spyMetaFetcher {
	return &spyMetaFetcher{}
}

func (s *spyMetaFetcher) Meta(ctx context.Context) (map[string]*logcache.MetaInfo, error) {
	s.ctx = ctx
	return s.results, s.err
}

type httpFailer struct{}

func (h *httpFailer) Do(req *http.Request) (*http.Response, error) {
	return nil, errors.New("some-error")
}

package handlers_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"time"

	logcache "code.cloudfoundry.org/go-log-cache"
	"code.cloudfoundry.org/go-loggregator/rpc/loggregator_v2"
	"code.cloudfoundry.org/loggregator-tools/request-spinner/internal/handlers"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Starter", func() {
	var (
		spyReader *spyReader
		recorder  *httptest.ResponseRecorder
		s         http.Handler
	)

	BeforeEach(func() {
		spyReader = newSpyReader()
		recorder = httptest.NewRecorder()
		s = handlers.NewStarter(spyReader)
	})

	It("reads from the given source ID", func() {
		req := httptest.NewRequest("POST", "/v1/start/?source_id=some-id", nil)
		s.ServeHTTP(recorder, req)

		Expect(recorder.Code).To(Equal(http.StatusOK))
		Expect(spyReader.sourceIDs).To(HaveLen(10))
		Expect(spyReader.sourceIDs[0]).To(Equal("some-id"))

		start := time.Now()
		Expect(spyReader.starts[0].UnixNano()).To(BeNumerically("~", time.Now().Add(-time.Minute).UnixNano(), time.Second))

		// Assert agains the delay
		Expect(time.Since(start)).To(BeNumerically("~", 10*time.Millisecond, 50*time.Millisecond))
	})

	It("uses the request's context", func() {
		ctx, cancel := context.WithCancel(context.Background())
		req := httptest.NewRequest("POST", "/v1/start?source_id=some-id", nil)
		req = req.WithContext(ctx)
		s.ServeHTTP(recorder, req)

		cancel()
		Expect(spyReader.ctxs[0].Done()).To(BeClosed())
	})

	It("stops reading upon an error", func() {
		spyReader.err = errors.New("some-error")

		req := httptest.NewRequest("POST", "/v1/start?source_id=some-id", nil)
		s.ServeHTTP(recorder, req)

		Expect(recorder.Code).To(Equal(http.StatusInternalServerError))
		Expect(spyReader.sourceIDs).To(HaveLen(1))
		Expect(recorder.Body.String()).To(Equal(`{"error": "some-error"}`))
	})

	It("reads with the given start time and cycles", func() {
		req := httptest.NewRequest("POST", "/v1/start?source_id=some-id&start_time=12345&cycles=11", nil)
		s.ServeHTTP(recorder, req)

		Expect(recorder.Code).To(Equal(http.StatusOK))
		Expect(spyReader.sourceIDs).To(HaveLen(11))
		Expect(spyReader.sourceIDs[0]).To(Equal("some-id"))
		Expect(spyReader.starts[0].UnixNano()).To(Equal(int64(12345)))
	})

	It("returns a 400 for requests with an invalid cycles", func() {
		req := httptest.NewRequest("POST", "/v1/start?source_id=some-id&cycles=-11", nil)
		s.ServeHTTP(recorder, req)

		Expect(recorder.Code).To(Equal(http.StatusBadRequest))
		Expect(recorder.Body.String()).To(Equal(`{"error": "failed to parse 'cycles' to an uint64"}`))
	})

	It("returns a 400 for requests with an invalid delay", func() {
		req := httptest.NewRequest("POST", "/v1/start?source_id=some-id&delay=invalid", nil)
		s.ServeHTTP(recorder, req)

		Expect(recorder.Code).To(Equal(http.StatusBadRequest))
		Expect(recorder.Body.String()).To(Equal(`{"error": "failed to parse 'delay' to a time.Duration"}`))
	})

	It("returns a 400 for requests with an invalid start_time", func() {
		req := httptest.NewRequest("POST", "/v1/start?source_id=some-id&start_time=invalid", nil)
		s.ServeHTTP(recorder, req)

		Expect(recorder.Code).To(Equal(http.StatusBadRequest))
		Expect(recorder.Body.String()).To(Equal(`{"error": "failed to parse 'start_time' to an int64"}`))
	})

	It("returns a 400 for requests without source ID", func() {
		req := httptest.NewRequest("POST", "/v1/start", nil)
		s.ServeHTTP(recorder, req)

		Expect(recorder.Code).To(Equal(http.StatusBadRequest))
		Expect(recorder.Body.String()).To(Equal(`{"error": "missing 'source_id' query parameter"}`))
	})

	It("returns a 404 for an unknown path", func() {
		req := httptest.NewRequest("GET", "/v1/unknown", nil)
		s.ServeHTTP(recorder, req)

		Expect(recorder.Code).To(Equal(http.StatusNotFound))
	})

	It("returns a 405 for a non-POST request", func() {
		req := httptest.NewRequest("GET", "/v1/start", nil)
		s.ServeHTTP(recorder, req)

		Expect(recorder.Code).To(Equal(http.StatusMethodNotAllowed))
	})
})

type spyReader struct {
	sourceIDs []string
	starts    []time.Time
	ctxs      []context.Context
	err       error
}

func newSpyReader() *spyReader {
	return &spyReader{}
}

func (s *spyReader) Read(ctx context.Context, sourceID string, start time.Time, opts ...logcache.ReadOption) ([]*loggregator_v2.Envelope, error) {
	s.sourceIDs = append(s.sourceIDs, sourceID)
	s.starts = append(s.starts, start)
	s.ctxs = append(s.ctxs, ctx)

	return nil, s.err
}

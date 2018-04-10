package expvarfilter_test

import (
	"fmt"

	"code.cloudfoundry.org/loggregator-tools/log-cache-forwarders/pkg/expvarfilter"

	"net/http"
	"net/http/httptest"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Handler", func() {
	var (
		rec *httptest.ResponseRecorder
		req *http.Request
	)

	BeforeEach(func() {
		var err error
		req, err = http.NewRequest("GET", "/", nil)
		Expect(err).ToNot(HaveOccurred())
		rec = httptest.NewRecorder()
	})

	It("calls the wrapped handler", func() {
		h := newSpyHandler()
		h.response = "{}"
		fh := expvarfilter.NewHandler(h, nil)

		fh.ServeHTTP(rec, req)

		Expect(h.called).To(BeTrue())
	})

	It("returns the response from the inner store", func() {
		h := newSpyHandler()
		h.response = "hello world"
		fh := expvarfilter.NewHandler(h, nil)

		fh.ServeHTTP(rec, req)

		Expect(rec.Body.Bytes()).To(BeEquivalentTo("hello world"))
	})

	It("only returns keys in the JSON response that are white listed", func() {
		h := newSpyHandler()
		h.response = `{
			"keyA": "valueA",
			"keyB": "valueB"
		}`
		fh := expvarfilter.NewHandler(h, []string{"keyA"})

		fh.ServeHTTP(rec, req)

		Expect(rec.Body.Bytes()).To(MatchJSON(`{"keyA": "valueA"}`))
	})
})

type spyHandler struct {
	called   bool
	response string
}

func newSpyHandler() *spyHandler {
	return &spyHandler{}
}

func (s *spyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.called = true
	fmt.Fprint(w, s.response)
}

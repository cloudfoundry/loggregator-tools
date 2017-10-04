package web_test

import (
	"net/http"
	"net/http/httptest"
	"strings"

	"code.cloudfoundry.org/loggregator-tools/syslog_to_datadog/internal/web"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Handler", func() {
	It("returns http status code Accepted", func() {
		var data []byte
		setter := func(d []byte) {
			data = d
		}

		h := web.NewHandler(setter)
		request := httptest.NewRequest("POST", "/", strings.NewReader("hello"))
		recorder := httptest.NewRecorder()

		h.ServeHTTP(recorder, request)
		Expect(recorder.Code).To(Equal(http.StatusAccepted))
		Expect(string(data)).To(Equal("hello"))
	})
})

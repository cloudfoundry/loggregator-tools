package web_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"

	"code.cloudfoundry.org/loggregator-tools/counter/web"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Counter", func() {
	It("sets the prime count", func() {
		body := strings.NewReader(`[{
			"id": "my-id",
			"primeCount": 1,
			"msgCount": 0
		}]`)
		req := httptest.NewRequest(http.MethodPut, "/set", body)
		c := web.NewCounter(10)
		rec := httptest.NewRecorder()

		c.SetHandler(rec, req)
		Expect(rec.Code).To(Equal(http.StatusOK))

		req = httptest.NewRequest(http.MethodGet, "/get-prime/my-id", nil)
		rec = httptest.NewRecorder()

		c.GetPrimeHandler(rec, req)
		Expect(rec.Code).To(Equal(http.StatusOK))
		Expect(rec.Body.String()).To(Equal("1"))
	})

	It("sets the message count", func() {
		body := strings.NewReader(`[{
			"id": "my-id",
			"primeCount": 0,
			"msgCount": 2
		}]`)
		req := httptest.NewRequest(http.MethodPut, "/set", body)
		c := web.NewCounter(10)
		rec := httptest.NewRecorder()

		c.SetHandler(rec, req)
		Expect(rec.Code).To(Equal(http.StatusOK))

		req = httptest.NewRequest(http.MethodGet, "/get/my-id", nil)
		rec = httptest.NewRecorder()

		c.GetHandler(rec, req)
		Expect(rec.Code).To(Equal(http.StatusOK))
		Expect(rec.Body.String()).To(Equal("2"))
	})

	It("returns zeros for unknown id", func() {
		req := httptest.NewRequest(http.MethodGet, "/get/my-unknown-id", nil)
		rec := httptest.NewRecorder()

		c := web.NewCounter(10)
		c.GetHandler(rec, req)
		Expect(rec.Code).To(Equal(http.StatusOK))
		Expect(rec.Body.String()).To(Equal("0"))

		req = httptest.NewRequest(http.MethodGet, "/get-prime/my-unknown-id", nil)
		rec = httptest.NewRecorder()

		c.GetPrimeHandler(rec, req)
		Expect(rec.Code).To(Equal(http.StatusOK))
		Expect(rec.Body.String()).To(Equal("0"))
	})

	It("keeps a limited number of counters", func() {
		c := web.NewCounter(2)

		for i := 0; i < 5; i++ {
			body := strings.NewReader(fmt.Sprintf(`[{
				"id": "my-id-%d",
				"primeCount": 100,
				"msgCount": 100}]`, i))
			req := httptest.NewRequest(http.MethodPut, "/set", body)
			rec := httptest.NewRecorder()

			c.SetHandler(rec, req)
		}

		for i := 0; i < 3; i++ {
			req := httptest.NewRequest(
				http.MethodGet,
				fmt.Sprintf("/get/my-id-%d", i),
				nil,
			)
			rec := httptest.NewRecorder()

			c.GetHandler(rec, req)
			Expect(rec.Body.String()).To(Equal("0"))
		}

		for i := 3; i < 5; i++ {
			req := httptest.NewRequest(
				http.MethodGet,
				fmt.Sprintf("/get/my-id-%d", i),
				nil,
			)
			rec := httptest.NewRecorder()

			c.GetHandler(rec, req)
			Expect(rec.Body.String()).To(Equal("100"))
		}
	})
})

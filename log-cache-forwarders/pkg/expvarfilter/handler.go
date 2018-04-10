package expvarfilter

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
)

type Handler struct {
	inner     http.Handler
	whitelist map[string]struct{}
}

type writeInterceptor struct {
	http.ResponseWriter
	b bytes.Buffer
}

func newWriteInterceptor(w http.ResponseWriter) *writeInterceptor {
	return &writeInterceptor{
		ResponseWriter: w,
	}
}

func (w *writeInterceptor) Write(b []byte) (int, error) {
	return w.b.Write(b)
}

func NewHandler(inner http.Handler, whitelist []string) Handler {
	w := make(map[string]struct{}, len(whitelist))
	for _, v := range whitelist {
		w[v] = struct{}{}
	}

	return Handler{
		inner:     inner,
		whitelist: w,
	}
}

func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	i := newWriteInterceptor(w)
	h.inner.ServeHTTP(i, r)

	var response map[string]interface{}
	err := json.Unmarshal(i.b.Bytes(), &response)
	if err != nil {
		panicWrite(w, i.b.Bytes())
		return
	}

	for key := range response {
		if _, ok := h.whitelist[key]; !ok {
			delete(response, key)
		}
	}

	resp, err := json.Marshal(response)
	if err != nil {
		log.Panic(err)
	}

	panicWrite(w, resp)
}

func panicWrite(w io.Writer, data []byte) {
	_, err := w.Write(data)
	if err != nil {
		log.Panic(err)
	}
}

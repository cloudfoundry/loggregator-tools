package web

import (
	"bytes"
	"log"
	"net/http"
)

// Setter is a func that can receive a slice of bytes.
type Setter func([]byte)

// Handler satisfies the http.Handler interface for receiving rfc 5424
// messages via HTTP.
type Handler struct {
	setter Setter
}

// NewHandler returns a new Handler.
func NewHandler(s Setter) *Handler {
	return &Handler{setter: s}
}

// ServeHTTP receives HTTP requests and calls the handlers Setter func with
// the request body.
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	buf := bytes.NewBuffer(make([]byte, 0, r.ContentLength))
	_, err := buf.ReadFrom(r.Body)
	r.Body.Close()
	if err != nil {
		log.Fatalf("failed to read request body %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	h.setter(buf.Bytes())

	w.WriteHeader(http.StatusAccepted)
}

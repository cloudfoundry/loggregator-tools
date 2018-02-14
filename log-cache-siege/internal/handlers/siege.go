package handlers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/url"

	"code.cloudfoundry.org/go-log-cache/rpc/logcache"
)

type Siege struct {
	f MetaFetcher
	c HTTPClient
	u url.URL
}

type MetaFetcher interface {
	Meta(ctx context.Context) (map[string]*logcache.MetaInfo, error)
}

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

func NewSiege(addr string, c HTTPClient, f MetaFetcher) *Siege {
	u, err := url.Parse(addr)
	if err != nil {
		log.Panicf("failed to parse URL: %s", err)
	}

	return &Siege{
		f: f,
		c: c,

		// Save a copy of the URL to ensure we can alter it without races.
		u: *u,
	}
}

func (s *Siege) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/v1/siege" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	m, err := s.f.Meta(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf(`{"error":%q}`, err)))
		return
	}

	for sourceID := range m {
		s.u.Path = fmt.Sprintf("/v1/start/%s", sourceID)
		req, err := http.NewRequest("POST", s.u.String(), nil)
		if err != nil {
			log.Panicf("failed to create request: %s", err)
		}

		resp, err := s.c.Do(req)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf(`{"error":%q}`, err)))
			return
		}

		if resp.StatusCode != http.StatusOK {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf(`{"error": "unexpected status code from request spinner: %d"}`, resp.StatusCode)))
		}
	}
}

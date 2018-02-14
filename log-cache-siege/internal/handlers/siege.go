package handlers

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"code.cloudfoundry.org/go-log-cache/rpc/logcache"
)

type Siege struct {
	f                  MetaFetcher
	c                  HTTPClient
	requestSpinnerAddr string
}

type MetaFetcher interface {
	Meta(ctx context.Context) (map[string]*logcache.MetaInfo, error)
}

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

func NewSiege(addr string, c HTTPClient, f MetaFetcher) *Siege {
	return &Siege{
		f:                  f,
		c:                  c,
		requestSpinnerAddr: addr,
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
		addr := fmt.Sprintf("%s/v1/start?source_id=%s", s.requestSpinnerAddr, sourceID)

		req, err := http.NewRequest("POST", addr, nil)
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

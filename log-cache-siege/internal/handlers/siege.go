package handlers

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"sync"

	"code.cloudfoundry.org/go-log-cache/rpc/logcache_v1"
)

type Siege struct {
	f                  MetaFetcher
	c                  HTTPClient
	requestSpinnerAddr string
	concurrentRequests int
}

type MetaFetcher interface {
	Meta(ctx context.Context) (map[string]*logcache_v1.MetaInfo, error)
}

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

func NewSiege(addr string, concurrentRequests int, c HTTPClient, f MetaFetcher) *Siege {
	return &Siege{
		f:                  f,
		c:                  c,
		requestSpinnerAddr: addr,
		concurrentRequests: concurrentRequests,
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

	sourceIDs := make(chan string)

	wg := sync.WaitGroup{}
	defer wg.Wait()
	for i := 0; i < s.concurrentRequests; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for sourceID := range sourceIDs {
				s.sendSpinnerRequest(sourceID)
			}
		}()
	}

	for sourceID := range m {
		sourceIDs <- sourceID
	}
	close(sourceIDs)
}

func (s *Siege) sendSpinnerRequest(sourceID string) {
	addr := fmt.Sprintf("%s/v1/start?source_id=%s", s.requestSpinnerAddr, sourceID)

	req, err := http.NewRequest("POST", addr, nil)
	if err != nil {
		log.Panicf("failed to create request: %s", err)
	}

	resp, err := s.c.Do(req)
	if err != nil {
		log.Printf("%s request to request spinner failed: %s", sourceID, err)
		return
	}
	defer func() {
		io.Copy(ioutil.Discard, resp.Body)
		resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		data, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Printf("failed to read request-spinner body: %s", err)
			return
		}

		log.Printf("%s unexpected status code from request spinner: %d -> %s", sourceID, resp.StatusCode, data)
		return
	}
}

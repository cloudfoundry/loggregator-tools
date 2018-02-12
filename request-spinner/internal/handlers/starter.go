package handlers

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	logcache "code.cloudfoundry.org/go-log-cache"
	"code.cloudfoundry.org/go-loggregator/rpc/loggregator_v2"
)

// Reader reads envelopes from Log Cache
type Reader interface {
	Read(
		ctx context.Context,
		sourceID string,
		start time.Time,
		opts ...logcache.ReadOption,
	) ([]*loggregator_v2.Envelope, error)
}

// Starter handles HTTP requests to start the request-spinner
type Starter struct {
	client Reader
}

// NewStarter creates a Starter
func NewStarter(r Reader) Starter {
	return Starter{
		client: r,
	}
}

// Starter implements http.Handler
func (s Starter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/v1/start" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	sourceID := r.URL.Query().Get("source_id")
	if sourceID == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": "missing 'source_id' query parameter"}`))
		return
	}

	startTime, err := startTime(r.URL.Query().Get("start_time"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": "failed to parse 'start_time' to an int64"}`))
		return
	}

	cycles, err := cycles(r.URL.Query().Get("cycles"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": "failed to parse 'cycles' to an uint64"}`))
		return
	}

	delay, err := delay(r.URL.Query().Get("delay"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": "failed to parse 'delay' to a time.Duration"}`))
		return
	}

	for i := 0; i < cycles; i++ {
		_, err := s.client.Read(r.Context(), sourceID, time.Unix(0, startTime))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf(`{"error": "%s"}`, err)))
			return
		}
		time.Sleep(delay)
	}
}

func startTime(param string) (int64, error) {
	if param == "" {
		return time.Now().Add(-time.Minute).UnixNano(), nil
	}

	return strconv.ParseInt(param, 10, 64)
}

func cycles(param string) (int, error) {
	if param == "" {
		return 10, nil
	}

	u, err := strconv.ParseUint(param, 10, 64)
	if err != nil {
		return 0, err
	}

	return int(u), nil
}

func delay(param string) (time.Duration, error) {
	if param == "" {
		return time.Millisecond, nil
	}

	return time.ParseDuration(param)
}

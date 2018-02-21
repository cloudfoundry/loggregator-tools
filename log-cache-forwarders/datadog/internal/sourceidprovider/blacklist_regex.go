package sourceidprovider

import (
	"context"
	"log"
	"regexp"
	"time"

	rpc "code.cloudfoundry.org/go-log-cache/rpc/logcache_v1"
)

// BlacklistRegex
type BlacklistRegex struct {
	f MetaFetcher
	r *regexp.Regexp
}

// MetaFetcher returns meta information from LogCache.
type MetaFetcher interface {

	// Meta returns meta information from LogCache.
	Meta(ctx context.Context) (map[string]*rpc.MetaInfo, error)
}

// NewBlacklistRegex compiles the configured regex pattern. If the pattern
// fails, it will panic.
func NewBlacklistRegex(pattern string, f MetaFetcher) *BlacklistRegex {
	return &BlacklistRegex{
		f: f,
		r: regexp.MustCompile(pattern),
	}
}

// SourceIDs returns each source ID provided by the MetaFetcher that did NOT
// match the given regex pattern.
func (r *BlacklistRegex) SourceIDs() []string {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	mi, err := r.f.Meta(ctx)
	if err != nil {
		log.Printf("failed to read Meta information: %s", err)
		return nil
	}

	var results []string
	for sourceID := range mi {
		if r.r.MatchString(sourceID) {
			continue
		}

		results = append(results, sourceID)
	}

	return results
}

package sourceidprovider

import (
	"context"
	"log"
	"regexp"
	"time"

	rpc "code.cloudfoundry.org/go-log-cache/rpc/logcache_v1"
)

// Regex provides source IDs that are filtered by a regex.
type Regex struct {
	f           MetaFetcher
	isBlacklist bool
	r           *regexp.Regexp
}

// MetaFetcher returns meta information from LogCache.
type MetaFetcher interface {
	// Meta returns meta information from LogCache.
	Meta(ctx context.Context) (map[string]*rpc.MetaInfo, error)
}

// NewRegex compiles the configured regex pattern. If the pattern
// fails, it will panic.
func NewRegex(isBlacklist bool, pattern string, f MetaFetcher) *Regex {
	return &Regex{
		f:           f,
		isBlacklist: isBlacklist,
		r:           regexp.MustCompile(pattern),
	}
}

// SourceIDs returns each source ID provided by the MetaFetcher that did NOT
// match the given regex pattern.
func (r *Regex) SourceIDs() []string {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	mi, err := r.f.Meta(ctx)
	if err != nil {
		log.Printf("failed to read Meta information: %s", err)
		return nil
	}

	var results []string
	for sourceID := range mi {
		// blacklist match
		//    0        0    append
		//    1        0    continue
		//    0        1    continue
		//    1        1    append
		match := r.r.MatchString(sourceID)
		if r.isBlacklist && !match || !r.isBlacklist && match {
			results = append(results, sourceID)
		}
	}

	return results
}

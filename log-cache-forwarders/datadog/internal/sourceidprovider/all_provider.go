package sourceidprovider

import (
	"context"
	"log"
	"time"
)

type AllProvider struct {
	f MetaFetcher
}

func All(f MetaFetcher) *AllProvider {
	return &AllProvider{
		f: f,
	}
}

func (p *AllProvider) SourceIDs() []string {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	mi, err := p.f.Meta(ctx)
	if err != nil {
		log.Printf("failed to read Meta information: %s", err)
		return nil
	}

	var results []string
	for sourceID := range mi {
		if sourceID != "" {
			results = append(results, sourceID)
		}
	}

	return results
}

package rss

import (
	"context"
	"fmt"

	"github.com/mmcdole/gofeed"
)

// Fetcher encapsulates RSS fetching logic
type Fetcher struct {
	fp *gofeed.Parser
}

// NewFetcher creates a new RSS fetcher
func NewFetcher() *Fetcher {
	return &Fetcher{
		fp: gofeed.NewParser(),
	}
}

// Fetch retrieves and parses an RSS feed from the given URL
func (f *Fetcher) Fetch(ctx context.Context, url string) (*gofeed.Feed, error) {
	feed, err := f.fp.ParseURLWithContext(url, ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch feed %s: %w", url, err)
	}
	return feed, nil
}

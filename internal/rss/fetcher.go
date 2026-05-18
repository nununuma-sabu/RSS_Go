package rss

import (
	"context"
	"fmt"

	"github.com/mmcdole/gofeed"
)

// Fetcher はRSS取得のロジックをカプセル化します
type Fetcher struct {
	fp *gofeed.Parser
}

// NewFetcher は新しいRSSフェッチャーを作成します
func NewFetcher() *Fetcher {
	return &Fetcher{
		fp: gofeed.NewParser(),
	}
}

// Fetch は指定されたURLからRSSフィードを取得して解析します
func (f *Fetcher) Fetch(ctx context.Context, url string) (*gofeed.Feed, error) {
	feed, err := f.fp.ParseURLWithContext(url, ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch feed %s: %w", url, err)
	}
	return feed, nil
}

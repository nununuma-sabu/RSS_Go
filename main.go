package main

import (
	"context"
	"log"
	"time"

	"github.com/nununuma-sabu/RSS_Go/internal/config"
	"github.com/nununuma-sabu/RSS_Go/internal/rss"
)

func main() {
	log.Println("Starting RSS fetcher...")

	// 1. Load config
	cfg, err := config.LoadConfig("config.yaml")
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	// 2. Initialize fetcher
	fetcher := rss.NewFetcher()

	// 3. Fetch feeds
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	for _, feedConfig := range cfg.Feeds {
		log.Printf("Fetching feed: %s (%s)", feedConfig.URL, feedConfig.Category)
		
		feed, err := fetcher.Fetch(ctx, feedConfig.URL)
		if err != nil {
			log.Printf("Error fetching feed %s: %v", feedConfig.URL, err)
			continue
		}

		log.Printf("Successfully fetched: %s", feed.Title)
		log.Printf("Found %d items", len(feed.Items))

		// Print top 3 items for demonstration
		limit := 3
		if len(feed.Items) < limit {
			limit = len(feed.Items)
		}

		for i := 0; i < limit; i++ {
			item := feed.Items[i]
			log.Printf("  - [%s] %s (%s)", feedConfig.Category, item.Title, item.Link)
		}
	}
	
	log.Println("RSS fetching completed.")
}

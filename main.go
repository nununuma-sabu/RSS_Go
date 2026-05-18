package main

import (
	"context"
	"log"
	"time"

	"github.com/mmcdole/gofeed"
	"github.com/nununuma-sabu/RSS_Go/internal/config"
	"github.com/nununuma-sabu/RSS_Go/internal/rss"
	"github.com/nununuma-sabu/RSS_Go/internal/storage"
)

func main() {
	log.Println("Starting RSS fetcher...")

	// 1. Load config
	cfg, err := config.LoadConfig("config.yaml")
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	// 2. Initialize storage
	store, err := storage.NewStorage("fetched_articles.json")
	if err != nil {
		log.Fatalf("Error loading storage: %v", err)
	}

	// 3. Initialize fetcher
	fetcher := rss.NewFetcher()

	// 4. Fetch feeds
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var newItems []*gofeed.Item

	for _, feedConfig := range cfg.Feeds {
		log.Printf("Fetching feed: %s (%s)", feedConfig.URL, feedConfig.Category)
		
		feed, err := fetcher.Fetch(ctx, feedConfig.URL)
		if err != nil {
			log.Printf("Error fetching feed %s: %v", feedConfig.URL, err)
			continue
		}

		log.Printf("Successfully fetched: %s", feed.Title)
		log.Printf("Found %d total items", len(feed.Items))

		// Filter new items
		var feedNewItems []*gofeed.Item
		for _, item := range feed.Items {
			if !store.IsFetched(item.Link) {
				feedNewItems = append(feedNewItems, item)
				store.MarkFetched(item.Link)
				newItems = append(newItems, item)
			}
		}
		
		log.Printf("Found %d NEW items", len(feedNewItems))

		// Print top 3 new items for demonstration
		limit := 3
		if len(feedNewItems) < limit {
			limit = len(feedNewItems)
		}

		for i := 0; i < limit; i++ {
			item := feedNewItems[i]
			log.Printf("  - [NEW] [%s] %s (%s)", feedConfig.Category, item.Title, item.Link)
		}
	}

	// 5. Save state
	if len(newItems) > 0 {
		log.Printf("Saving state... Added %d new items.", len(newItems))
		if err := store.Save(); err != nil {
			log.Fatalf("Error saving storage: %v", err)
		}
	} else {
		log.Println("No new items to save.")
	}
	
	log.Println("RSS fetching completed.")
}

package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/mmcdole/gofeed"
	"github.com/nununuma-sabu/RSS_Go/internal/config"
	"github.com/nununuma-sabu/RSS_Go/internal/rss"
	"github.com/nununuma-sabu/RSS_Go/internal/storage"
	"github.com/nununuma-sabu/RSS_Go/internal/summarizer"
)

func main() {
	log.Println("Starting RSS fetcher and summarizer...")

	// Get API Key (Fail fast if not set)
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		log.Fatalf("Error: GEMINI_API_KEY environment variable is not set. Please provide it to run the summarizer.")
	}

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

	// 4. Initialize summarizer
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	sumz, err := summarizer.NewSummarizer(ctx, apiKey)
	if err != nil {
		log.Fatalf("Error initializing summarizer: %v", err) // API key value is NOT logged
	}
	defer sumz.Close()

	var newItems []*gofeed.Item

	// Ensure summaries directory exists
	if err := os.MkdirAll("summaries", 0755); err != nil {
		log.Fatalf("Error creating summaries directory: %v", err)
	}

	// Prepare output file (e.g. summaries/2026-05-18.md)
	today := time.Now().Format("2006-01-02")
	summaryFilename := fmt.Sprintf("summaries/%s.md", today)
	summaryFile, err := os.OpenFile(summaryFilename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Error opening summary output file: %v", err)
	}
	defer summaryFile.Close()

	info, err := summaryFile.Stat()
	if err == nil && info.Size() == 0 {
		summaryFile.WriteString(fmt.Sprintf("# %s の新着記事サマリー\n\n", today))
	}

	// 5. Fetch feeds and summarize
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

		// Summarize up to 3 items per feed to avoid rate limits / high latency
		limit := 3
		if len(feedNewItems) < limit {
			limit = len(feedNewItems)
		}

		for i := 0; i < limit; i++ {
			item := feedNewItems[i]
			log.Printf("  - Summarizing: [%s] %s", feedConfig.Category, item.Title)

			// Prepare text
			contentToSummarize := fmt.Sprintf("Title: %s\nLink: %s\nDescription: %s", item.Title, item.Link, item.Description)
			
			summary, err := sumz.Summarize(ctx, contentToSummarize)
			if err != nil {
				log.Printf("    Error summarizing %s: %v", item.Title, err)
				continue
			}

			// Write result to summary.md
			timestamp := time.Now().Format("2006/01/02 15:04:05")
			entry := fmt.Sprintf("## [%s] %s\n- **Date:** %s\n- **Link:** %s\n\n### 要約\n%s\n\n---\n", feedConfig.Category, item.Title, timestamp, item.Link, summary)
			if _, err := summaryFile.WriteString(entry); err != nil {
				log.Printf("    Error writing to summary.md: %v", err)
			}
			log.Printf("    -> Summarized and saved to summary.md")
		}
	}

	// 6. Save state
	if len(newItems) > 0 {
		log.Printf("Saving state... Added %d new items.", len(newItems))
		if err := store.Save(); err != nil {
			log.Fatalf("Error saving storage: %v", err)
		}
	} else {
		log.Println("No new items to save.")
		summaryFile.WriteString("本日の新着記事はありませんでした。\n\n")
	}
	
	log.Println("RSS fetching and summarization completed.")
}

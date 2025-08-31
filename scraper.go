package main

import (
	"context"
	"database/sql"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/7-Dany/dev/internal/database"
	"github.com/google/uuid"
)

func startScraping(db *database.Queries, concurrency int, timeBetweenRequests time.Duration) {
	log.Printf("Scraping on %v goroutines, every %v duration", concurrency, timeBetweenRequests)
	ticker := time.NewTicker(timeBetweenRequests)

	for ; ; <-ticker.C {
		feeds, err := db.GetNextFeedToFetch(
			context.Background(),
			int32(concurrency),
		)
		if err != nil {
			log.Printf("error fetching feeds, %v", err)
			continue
		}

		wg := &sync.WaitGroup{}
		for _, feed := range feeds {
			wg.Add(1)

			go scrapeFeed(db, feed, wg)
		}
		wg.Wait()
	}
}

func scrapeFeed(db *database.Queries, feed database.Feed, wg *sync.WaitGroup) {
	defer wg.Done()

	_, err := db.MarkFeedAsFetched(context.Background(), feed.ID)
	if err != nil {
		log.Printf("error marking feeds as fetched, %v", err)
		return
	}

	rssFeed, err := urlToFeed(feed.Url)
	if err != nil {
		log.Printf("error fetchin feed, %v", err)
		return
	}

	for _, item := range rssFeed.Channel.Item {
		description := sql.NullString{}
		if item.Description != "" {
			description.String = item.Description
			description.Valid = true
		}

		t, err := time.Parse(time.RFC1123Z, item.PubDate)
		if err != nil {
			log.Printf("couldn't parse item.PubDate: %v, %v", item.PubDate, err)
			continue
		}

		_, err = db.CreatePost(context.Background(), database.CreatePostParams{
			ID:          uuid.New(),
			FeedID:      feed.ID,
			Title:       item.Title,
			Description: description,
			Url:         item.Link,
			CreatedAt:   time.Now().UTC(),
			UpdatedAt:   time.Now().UTC(),
			PublishedAt: t,
		})
		if err != nil {
			if strings.Contains(err.Error(), "duplicate key") {
				continue
			}
			log.Printf("failed to create post, %v", err)
		}
	}

	log.Printf("Feed %s collected, %v posts found", feed.Name, len(rssFeed.Channel.Item))
}

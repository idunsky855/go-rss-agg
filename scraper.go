package main

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/idunsky855/go-rss-agg/internal/database"
)

func startScraping(db *database.Queries, concurrency int, timeBetweenRequests time.Duration) {
	log.Printf("scraping on %v goroutines every %s duration", concurrency, timeBetweenRequests)
	ticker := time.NewTicker(timeBetweenRequests)
	for ; ; <-ticker.C {
		feeds, err := db.GetNextFeedsToFetch(
			context.Background(),
			int32(concurrency),
		)
		if err != nil {
			log.Println("error fetching feeds: ", err)
			continue
		}

		wg := &sync.WaitGroup{}
		for _, feed := range feeds {
			wg.Add(1)

			go scrapeFeed(db, wg, feed)
		}
		wg.Wait()
	}

}

func scrapeFeed(db *database.Queries, wg *sync.WaitGroup, feed database.Feed) {
	defer wg.Done()
	_, err := db.MarkFeedAsFetched(context.Background(), feed.ID)
	if err != nil {
		log.Println("Error marking feed as fetched")
		return
	}

	rssFeed, err := URLToFeed(feed.Url)
	if err != nil {
		log.Println("Error fetching feed: ", err)
		return
	}

	for _, item := range rssFeed.Channel.Item {
		log.Println("Found post", item.Title, "on feed", feed.Name)
	}

	log.Println("Feed %s collected %v posts found", feed.Name, len(rssFeed.Channel.Item))
}

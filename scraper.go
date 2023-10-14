package main

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/witczakxd/first-web-api/internal/database"
)

func startScraping(
	db *database.Queries,
	concurrency int,
	timeBetweenRequest time.Duration,
) {
	log.Printf("Scraping on %v goroutines every %s duration", concurrency, timeBetweenRequest)
	ticker := time.NewTicker(timeBetweenRequest)
	for ; ; <-ticker.C {
		feeds,err := db.GetNextFeedsToFetch(
			context.Background(),
			int32(concurrency),
		)
		if err != nil {
			log.Println("Error while fetching feeds",err)
			continue
		}

		wg := &sync.WaitGroup{}
		for _,feed := range feeds {
			wg.Add(1)
			go scrapeFeed(db,wg,feed)
		}
		wg.Wait()
	}
}

func scrapeFeed (db *database.Queries,wg *sync.WaitGroup,feed database.Feed) {
	defer wg.Done()

	_,err := db.MarkFeedAsFetched(context.Background(),feed.ID)
	if err != nil {
		log.Println("Error while marking feed as fetched",err)
		return
	}

	rssFeed,err := urlToFeed(feed.Url)
	if err != nil {
		log.Println("Error while fetching feed",err)
		return
	}

	for _,item := range rssFeed.Channel.Item {
		log.Println("Found item",item.Title, "on feed ",feed.Name)
	}
	log.Printf("Scraped feed %s, %v posts found",feed.Name,len(rssFeed.Channel.Item))
}
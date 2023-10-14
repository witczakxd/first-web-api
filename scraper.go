package main

import (
	"context"
	"database/sql"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
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
		description := sql.NullString{}
		if item.Description != "" {
			description.Valid = true
			description.String = item.Description
		}

		pubAt,err := time.Parse(time.RFC1123Z,item.PubDate)
		if err != nil {
			log.Println("Error while parsing date",err)
			continue
		}

		_,err = db.CreatePost(context.Background(),database.CreatePostParams{
			ID:			uuid.New(),
			CreatedAt: 	time.Now().UTC(),
			UpdatedAt: 	time.Now().UTC(),
			Title: 		item.Title,
			Description: description,
			PublishedAt: pubAt,
			Url: 		item.Link,
			FeedID: 	feed.ID,
		})

		if err != nil {
			if strings.Contains(err.Error(),"unique constraint") {
				continue
			}
			log.Println("Error while creating post",err)
			continue
		}

	}
	log.Printf("Scraped feed %s, %v posts found",feed.Name,len(rssFeed.Channel.Item))
}
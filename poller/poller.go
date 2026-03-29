package poller

import (
	"github.com/mmcdole/gofeed"
	"gross/db"
	"gross/models"
	"log"
	"time"
)

func Start(interval time.Duration) {
	ticker := time.NewTicker(interval)

	go func() {
		pollAll()
		for range ticker.C {
			pollAll()
		}
	}()

	log.Printf("poller started - refreshing feeds every %s\n", interval)
}

func pollAll() {
	feeds, err := db.GetAllFeeds()
	if err != nil {
		log.Println("poller: failed to load feeds:", err)
		return
	}

	sem := make(chan struct{}, 5) // max 5 concurrent fetches
	for _, feed := range feeds {
		sem <- struct{}{}
		go func(f models.Feed) {
			defer func() { <-sem }()
			if err := pollFeed(f); err != nil {
				log.Printf("poller: error fetching %s: %v\n", f.URL, err)
			}
		}(feed)
	}

	// Wait for all workers to finish
	for i := 0; i < cap(sem); i++ {
		sem <- struct{}{}
	}
}

func pollFeed(f models.Feed) error {
	fp := gofeed.NewParser()
	feed, err := fp.ParseURL(f.URL)
	if err != nil {
		return err
	}

	newCount := 0
	for _, item := range feed.Items {
		// item.PublishedParsed can sometimes be nil — guard against it
		var published time.Time
		if item.PublishedParsed != nil {
			published = *item.PublishedParsed
		} else {
			published = time.Now()
		}

		err := db.InsertArticle(models.Article{
			FeedID:      f.ID,
			Title:       item.Title,
			Link:        item.Link,
			Description: item.Description,
			Published:   published,
		})
		if err == nil {
			newCount++
		}
	}

	db.UpdateLastFetched(f.ID)
	log.Printf("poller: %s — %d new articles\n", f.Title, newCount)
	return nil
}

package handlers

import (
	"github.com/mmcdole/gofeed"
	"log"
	_ "time"
)

func FetchFeed(url string) error {
	fp := gofeed.NewParser()
	feed, err := fp.ParseURL(url)
	if err != nil {
		return err
	}

	for _, item := range feed.Items {
		log.Printf("Title: %s | Published: %s\n", item.Title, item.Published)
		//TODO: save to DB, skip dupes by checking item.Link
	}

	return nil
}

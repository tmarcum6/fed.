package models

import "time"

type Feed struct {
	ID          int
	URL         string
	Title       string
	LastFetched time.Time
}

type Article struct {
	ID          int
	FeedID      int
	Title       string
	Link        string
	Description string
	Published   time.Time
	Read        bool
}

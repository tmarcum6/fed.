package poller

import (
	"gross/db"
	"gross/models"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPollFeed(t *testing.T) {
	db.Init(":memory:")

	// Spin up a fake RSS server so tests don't hit the real internet
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/rss+xml")
		w.Write([]byte(`<?xml version="1.0"?>
        <rss version="2.0">
          <channel>
            <title>Test Feed</title>
            <item>
              <title>Test Article</title>
              <link>https://example.com/article-1</link>
              <description>Hello world</description>
            </item>
          </channel>
        </rss>`))
	}))
	defer server.Close()

	feed := models.Feed{ID: 1, URL: server.URL, Title: "Test Feed"}
	db.InsertFeed(feed.URL, feed.Title)

	err := pollFeed(feed)
	if err != nil {
		t.Fatalf("pollFeed failed: %v", err)
	}

	articles, _ := db.GetArticles("1", false, false, false)
	if len(articles) == 0 {
		t.Fatal("expected articles to be saved after polling")
	}
	if articles[0].Title != "Test Article" {
		t.Errorf("expected 'Test Article', got '%s'", articles[0].Title)
	}
}

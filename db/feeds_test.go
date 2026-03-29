package db

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	// Use a temporary DB for tests so you don't pollute rss.db
	Init(":memory:") // SQLite in-memory DB — wiped after tests finish
	code := m.Run()
	os.Exit(code)
}

func TestInsertAndGetFeed(t *testing.T) {
	id, err := InsertFeed("https://example.com/feed.xml", "Example Blog")
	if err != nil {
		t.Fatalf("InsertFeed failed: %v", err)
	}
	if id == 0 {
		t.Fatal("expected non-zero ID")
	}

	feeds, err := GetAllFeeds()
	if err != nil {
		t.Fatalf("GetAllFeeds failed: %v", err)
	}
	if len(feeds) == 0 {
		t.Fatal("expected at least one feed")
	}
	if feeds[0].Title != "Example Blog" {
		t.Errorf("expected title 'Example Blog', got '%s'", feeds[0].Title)
	}
}

func TestInsertDuplicateFeed(t *testing.T) {
	InsertFeed("https://duplicate.com/feed.xml", "Dupe")
	_, err := InsertFeed("https://duplicate.com/feed.xml", "Dupe")
	// INSERT OR IGNORE means this should NOT error
	if err != nil {
		t.Errorf("duplicate insert should be silently ignored, got: %v", err)
	}
}

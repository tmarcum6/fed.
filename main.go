package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/mmcdole/gofeed"
	"gross/db"
	"gross/models"
	"gross/poller"
	"gross/templates"
	"log"
	"net/http"
	"strconv"
	"time"
)

type PageData struct {
	Feeds      []models.Feed
	Articles   []models.Article
	Feed       *models.Feed
	UnreadOnly bool
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	unreadOnly := r.URL.Query().Get("unread") == "true"
	feeds, _ := db.GetAllFeeds()
	articles, _ := db.GetArticles("", unreadOnly)

	templates.Render(w, "index.html", PageData{
		Feeds:      feeds,
		Articles:   articles,
		UnreadOnly: unreadOnly,
	})
}

func feedHandler(w http.ResponseWriter, r *http.Request) {
	feedID := chi.URLParam(r, "id")
	feeds, _ := db.GetAllFeeds()
	articles, _ := db.GetArticles(feedID, false)

	var current *models.Feed
	for _, f := range feeds {
		if strconv.Itoa(f.ID) == feedID {
			current = &f
			break
		}
	}

	templates.Render(w, "feed.html", PageData{
		Feeds:    feeds,
		Articles: articles,
		Feed:     current,
	})
}

func addFeedHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	url := r.FormValue("url")
	if url == "" {
		http.Error(w, "url is required", http.StatusBadRequest)
		return
	}

	fp := gofeed.NewParser()
	feed, err := fp.ParseURL(url)
	if err != nil {
		http.Error(w, "could not parse feed: "+err.Error(), http.StatusBadRequest)
		return
	}

	feedID, err := db.InsertFeed(url, feed.Title)
	if err != nil {
		http.Error(w, "failed to save feed", http.StatusInternalServerError)
		return
	}

	for _, item := range feed.Items {
		var published time.Time
		if item.PublishedParsed != nil {
			published = *item.PublishedParsed
		} else {
			published = time.Now()
		}
		db.InsertArticle(models.Article{
			FeedID:      int(feedID),
			Title:       item.Title,
			Link:        item.Link,
			Description: item.Description,
			Published:   published,
		})
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// PATCH /articles/{id}/read
func markReadHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	articleID, _ := strconv.Atoi(id)

	if err := db.MarkAsRead(articleID); err != nil {
		http.Error(w, "failed to mark as read", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func main() {
	db.Init("./rss.db")
	templates.Load()
	poller.Start(15 * time.Minute)

	r := chi.NewRouter()
	r.Get("/", indexHandler)
	r.Post("/feeds", addFeedHandler)
	r.Get("/feeds/{id}", feedHandler)
	r.Post("/articles/{id}/read", markReadHandler)

	log.Println("Listening on :8080")
	http.ListenAndServe(":8080", r)
}

package main

import (
	"encoding/json"
	"gross/db"
	"gross/handlers"
	"gross/models"
	"gross/poller"
	"gross/templates"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/mmcdole/gofeed"
)

type PageData struct {
	Feeds      []models.Feed
	Articles   []models.Article
	Feed       *models.Feed
	UnreadOnly bool
	Hidden     bool
}

type FeedStats struct {
	models.Feed
	Total  int
	Unread int
}

type FeedsPageData struct {
	Feeds []FeedStats
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	unreadOnly := r.URL.Query().Get("unread") == "true"
	hiddenOnly := r.URL.Query().Get("hidden") == "true"
	feeds, _ := db.GetAllFeeds()
	articles, _ := db.GetArticles("", unreadOnly, hiddenOnly, false)

	templates.Render(w, "index.html", PageData{
		Feeds:      feeds,
		Articles:   articles,
		UnreadOnly: unreadOnly,
		Hidden:     hiddenOnly,
	})
}

func feedHandler(w http.ResponseWriter, r *http.Request) {
	feedID := chi.URLParam(r, "id")
	unreadOnly := r.URL.Query().Get("unread") == "true"
	feeds, _ := db.GetAllFeeds()
	articles, _ := db.GetArticles(feedID, unreadOnly, false, false)

	var current *models.Feed
	for _, f := range feeds {
		if strconv.Itoa(f.ID) == feedID {
			current = &f
			break
		}
	}

	templates.Render(w, "feed.html", PageData{
		Feeds:      feeds,
		Articles:   articles,
		Feed:       current,
		UnreadOnly: unreadOnly,
	})
}

func savedHandler(w http.ResponseWriter, r *http.Request) {
	feeds, _ := db.GetAllFeeds()
	articles, _ := db.GetArticles("", false, false, true)

	templates.Render(w, "saved.html", PageData{
		Feeds:    feeds,
		Articles: articles,
	})
}

func feedsHandler(w http.ResponseWriter, r *http.Request) {
	feeds, _ := db.GetAllFeeds()
	var feedsWithStats []FeedStats
	for _, f := range feeds {
		total, unread, _ := db.GetFeedStats(f.ID)
		feedsWithStats = append(feedsWithStats, FeedStats{
			Feed:   f,
			Total:  total,
			Unread: unread,
		})
	}

	templates.Render(w, "feeds.html", FeedsPageData{
		Feeds: feedsWithStats,
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

	if feedID > 0 {
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
	}

	http.Redirect(w, r, "/feeds", http.StatusSeeOther)
}

func refreshFeedHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	feedID, _ := strconv.Atoi(id)

	feed, err := db.GetFeedByID(feedID)
	if err != nil {
		http.Error(w, "feed not found", http.StatusNotFound)
		return
	}

	fp := gofeed.NewParser()
	fetchedFeed, err := fp.ParseURL(feed.URL)
	if err != nil {
		http.Error(w, "failed to fetch feed: "+err.Error(), http.StatusBadRequest)
		return
	}

	for _, item := range fetchedFeed.Items {
		var published time.Time
		if item.PublishedParsed != nil {
			published = *item.PublishedParsed
		} else {
			published = time.Now()
		}
		db.InsertArticle(models.Article{
			FeedID:      feedID,
			Title:       item.Title,
			Link:        item.Link,
			Description: item.Description,
			Published:   published,
		})
	}

	db.UpdateLastFetched(feedID)

	http.Redirect(w, r, "/feeds", http.StatusSeeOther)
}

func updateFeedHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	feedID, _ := strconv.Atoi(id)

	var req struct {
		URL   string `json:"url"`
		Title string `json:"title"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	if req.URL == "" {
		http.Error(w, "url is required", http.StatusBadRequest)
		return
	}

	if req.Title == "" {
		fp := gofeed.NewParser()
		feed, err := fp.ParseURL(req.URL)
		if err != nil {
			http.Error(w, "could not parse feed: "+err.Error(), http.StatusBadRequest)
			return
		}
		req.Title = feed.Title
	}

	if err := db.UpdateFeedURL(feedID, req.URL, req.Title); err != nil {
		http.Error(w, "failed to update feed", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func deleteFeedHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	feedID, _ := strconv.Atoi(id)

	db.DeleteArticlesByFeed(feedID)
	db.DeleteFeedByID(feedID)

	w.WriteHeader(http.StatusNoContent)
}

func markReadHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	articleID, _ := strconv.Atoi(id)

	if err := db.MarkAsRead(articleID); err != nil {
		http.Error(w, "failed to mark as read", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func markHiddenHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	articleID, _ := strconv.Atoi(id)

	if err := db.MarkAsHidden(articleID); err != nil {
		http.Error(w, "failed to mark as hidden", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func unhideArticleHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	articleID, _ := strconv.Atoi(id)

	if err := db.UnhideArticle(articleID); err != nil {
		http.Error(w, "failed to unhide article", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func markUnreadHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	articleID, _ := strconv.Atoi(id)

	if err := db.MarkAsUnread(articleID); err != nil {
		http.Error(w, "failed to mark as unread", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func saveArticleHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	articleID, _ := strconv.Atoi(id)

	if err := db.SaveArticle(articleID); err != nil {
		http.Error(w, "failed to save article", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func unsaveArticleHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	articleID, _ := strconv.Atoi(id)

	if err := db.UnsaveArticle(articleID); err != nil {
		http.Error(w, "failed to unsave article", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func exportFeedsHandler(w http.ResponseWriter, r *http.Request) {
	feeds, err := db.GetAllFeeds()
	if err != nil {
		http.Error(w, "failed to get feeds", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/csv")
	w.Header().Set("Content-Disposition", "attachment; filename=feeds.csv")

	for _, feed := range feeds {
		w.Write([]byte(feed.URL + "\n"))
	}
}

func backupHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Disposition", "attachment; filename=fed.db")
	http.ServeFile(w, r, "./rss.db")
}

func main() {
	db.Init("./rss.db")
	templates.Load()
	poller.Start(15 * time.Minute)

	r := chi.NewRouter()
	r.Get("/", indexHandler)
	r.Get("/saved", savedHandler)
	r.Get("/feeds", feedsHandler)
	r.Post("/feeds", addFeedHandler)
	r.Post("/feeds/import", handlers.ImportFeedsCSVHandler)
	r.Get("/api/export-feeds", exportFeedsHandler)
	r.Get("/feeds/discover", handlers.DiscoverFeedsHandler)
	r.Post("/feeds/from-discovery", handlers.AddFeedFromDiscoveryHandler)
	r.Post("/feeds/{id}/refresh", refreshFeedHandler)
	r.Put("/feeds/{id}", updateFeedHandler)
	r.Get("/backup", backupHandler)
	r.Get("/feeds/{id}", feedHandler)
	r.Post("/articles/{id}/read", markReadHandler)
	r.Post("/articles/{id}/unread", markUnreadHandler)
	r.Post("/articles/{id}/hidden", markHiddenHandler)
	r.Post("/articles/{id}/unhide", unhideArticleHandler)
	r.Post("/articles/{id}/save", saveArticleHandler)
	r.Post("/articles/{id}/unsave", unsaveArticleHandler)
	r.Delete("/feeds/{id}", deleteFeedHandler)

	log.Println("Listening on :8080")
	http.ListenAndServe(":8080", r)
}

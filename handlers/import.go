package handlers

import (
	"encoding/json"
	"gross/db"
	"gross/models"
	"mime/multipart"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/mmcdole/gofeed"
)

const (
	maxFeedsPerImport = 100
	maxFileSize       = 1 << 20 // 1MB
	feedParseTimeout  = 10 * time.Second
)

type ImportResult struct {
	Total   int      `json:"total"`
	Added   int      `json:"added"`
	Skipped int      `json:"skipped"`
	Errors  []string `json:"errors"`
}

type ImportError struct {
	URL   string `json:"url"`
	Error string `json:"error"`
}

var (
	urlRegex    = regexp.MustCompile(`^https?://[^\s]+$`)
	importMutex sync.Mutex
)

func ImportFeedsCSVHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseMultipartForm(maxFileSize); err != nil {
		http.Error(w, "File too large or invalid form", http.StatusBadRequest)
		return
	}

	file, _, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "No file provided", http.StatusBadRequest)
		return
	}
	defer file.Close()

	urls, err := parseCSVFile(file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if len(urls) == 0 {
		sendJSONResponse(w, ImportResult{
			Total:   0,
			Added:   0,
			Skipped: 0,
			Errors:  []string{"No valid URLs found in file"},
		})
		return
	}

	result := importFeeds(urls)
	sendJSONResponse(w, result)
}

func parseCSVFile(file multipart.File) ([]string, error) {
	content := make([]byte, maxFileSize)
	n, err := file.Read(content)
	if err != nil && err.Error() != "EOF" {
		return nil, err
	}
	content = content[:n]

	raw := strings.ReplaceAll(string(content), "\r\n", "\n")
	raw = strings.ReplaceAll(raw, "\r", "\n")

	parts := strings.FieldsFunc(raw, func(r rune) bool {
		return r == ',' || r == '\n'
	})

	var urls []string
	seen := make(map[string]bool)

	for _, part := range parts {
		part = strings.TrimSpace(part)
		part = strings.Trim(part, "\"")

		if part == "" {
			continue
		}

		if seen[part] {
			continue
		}

		parsed, err := url.Parse(part)
		if err != nil {
			continue
		}

		if parsed.Scheme != "https" && parsed.Scheme != "http" {
			continue
		}

		if !urlRegex.MatchString(part) {
			continue
		}

		urls = append(urls, part)
		seen[part] = true

		if len(urls) >= maxFeedsPerImport {
			break
		}
	}

	return urls, nil
}

func importFeeds(inputURLs []string) ImportResult {
	importMutex.Lock()
	defer importMutex.Unlock()

	existingURLs, err := db.GetAllFeedURLs()
	if err != nil {
		return ImportResult{
			Total:  len(inputURLs),
			Errors: []string{"Failed to check existing feeds: " + err.Error()},
		}
	}

	existingSet := make(map[string]bool)
	for _, u := range existingURLs {
		existingSet[u] = true
	}

	var added, skipped int
	var errors []string

	for _, feedURL := range inputURLs {
		if existingSet[feedURL] {
			skipped++
			continue
		}

		fp := gofeed.NewParser()
		feed, err := fp.ParseURL(feedURL)
		if err != nil {
			errors = append(errors, err.Error()+" ("+feedURL+")")
			continue
		}

		title := feedURL
		if feed != nil && feed.Title != "" {
			title = feed.Title
		}

		feedID, err := db.InsertFeed(feedURL, title)
		if err != nil {
			errors = append(errors, "Failed to insert feed: "+feedURL)
			continue
		}

		if feedID > 0 && feed != nil {
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
			added++
		}
	}

	if len(errors) > 10 {
		errors = append(errors[:10], "...")
	}

	return ImportResult{
		Total:   len(inputURLs),
		Added:   added,
		Skipped: skipped,
		Errors:  errors,
	}
}

func sendJSONResponse(w http.ResponseWriter, result ImportResult) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
}

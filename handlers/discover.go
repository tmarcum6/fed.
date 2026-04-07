package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"github.com/mmcdole/gofeed"
)

type DiscoveredFeed struct {
	Title string `json:"title"`
	URL   string `json:"url"`
	Type  string `json:"type"`
}

type DiscoverResponse struct {
	Feeds     []DiscoveredFeed `json:"feeds"`
	DirectURL string           `json:"directUrl,omitempty"`
}

var commonFeedPaths = []string{
	"/feed",
	"/feed/",
	"/rss",
	"/rss.xml",
	"/feed.xml",
	"/atom.xml",
	"/index.xml",
	"/rss/feed.xml",
	"/blog/feed",
	"/posts/feed",
	"/feed/rss",
	"/feed/atom",
}

func DiscoverFeedsHandler(w http.ResponseWriter, r *http.Request) {
	siteURL := r.URL.Query().Get("url")
	if siteURL == "" {
		http.Error(w, "url parameter is required", http.StatusBadRequest)
		return
	}

	feeds, err := DiscoverFeeds(siteURL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(feeds)
}

func DiscoverFeeds(siteURL string) ([]DiscoveredFeed, error) {
	var feeds []DiscoveredFeed
	seen := make(map[string]bool)

	parsedURL, err := url.Parse(siteURL)
	if err != nil {
		return nil, err
	}

	baseURL := parsedURL.Scheme + "://" + parsedURL.Host
	if port := parsedURL.Port(); port != "" {
		baseURL += ":" + port
	}

	fp := gofeed.NewParser()

	directFeed, err := fp.ParseURL(siteURL)
	if err == nil && directFeed.Title != "" {
		feeds = append(feeds, DiscoveredFeed{
			Title: directFeed.Title,
			URL:   siteURL,
			Type:  "direct",
		})
		seen[siteURL] = true
	}

	client := &http.Client{}
	for _, path := range commonFeedPaths {
		feedURL := baseURL + path
		if seen[feedURL] {
			continue
		}

		resp, err := client.Head(feedURL)
		if err != nil {
			continue
		}
		resp.Body.Close()

		if resp.StatusCode == http.StatusOK {
			contentType := resp.Header.Get("Content-Type")
			if strings.Contains(contentType, "xml") || strings.Contains(contentType, "rss") || strings.Contains(contentType, "atom") {
				feed, err := fp.ParseURL(feedURL)
				title := feedURL
				if err == nil && feed.Title != "" {
					title = feed.Title
				}
				feeds = append(feeds, DiscoveredFeed{
					Title: title,
					URL:   feedURL,
					Type:  "discovered",
				})
				seen[feedURL] = true
			}
		}
	}

	htmlURL := baseURL + "/"
	resp, err := http.Get(htmlURL)
	if err == nil {
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err == nil {
			html := string(body)

			linkPattern := regexp.MustCompile(`<link[^>]+type=["']?(application/rss\+xml|application/atom\+xml)["']?[^>]+href=["']([^"']+)["'][^>]*>`)
			titlePattern := regexp.MustCompile(`<link[^>]+href=["'][^"']+["'][^>]+title=["']([^"']+)["'][^>]*>`)

			matches := linkPattern.FindAllStringSubmatch(html, -1)
			for _, match := range matches {
				feedURL := match[2]
				if !strings.HasPrefix(feedURL, "http") {
					if strings.HasPrefix(feedURL, "/") {
						feedURL = baseURL + feedURL
					} else {
						feedURL = baseURL + "/" + feedURL
					}
				}

				if seen[feedURL] {
					continue
				}

				title := feedURL
				titleMatch := titlePattern.FindStringSubmatch(html)
				if titleMatch != nil {
					title = titleMatch[1]
				}

				feedType := "rss"
				if strings.Contains(match[1], "atom") {
					feedType = "atom"
				}

				feeds = append(feeds, DiscoveredFeed{
					Title: title,
					URL:   feedURL,
					Type:  feedType,
				})
				seen[feedURL] = true
			}
		}
	}

	return feeds, nil
}

func AddFeedFromDiscoveryHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		URL string `json:"url"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	fp := gofeed.NewParser()
	_, err := fp.ParseURL(req.URL)
	if err != nil {
		http.Error(w, "could not parse feed: "+err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"url": req.URL,
	})
}

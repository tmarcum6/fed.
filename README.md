# fed.

**fed.** is a lightweight, self-hosted RSS feed reader built in Go. No Electron. No Node. No nonsense. Just a single binary, a SQLite database, and your feeds.

---

## Features

- **Add any RSS or Atom feed** — via URL
- **Discover feeds** — auto-detect RSS/Atom feeds on any website
- **Import feeds** — bulk import from CSV (one URL per line)
- **Auto-polling** — feeds refresh every 15 minutes in the background
- **Filters** — view all, unread only, or hidden articles
- **Per-feed filtering** — toggle unread filter on individual feeds
- **Mark as read/unread** — track what you've seen
- **Hide articles** — hide articles you don't want to see
- **Manage feeds** — edit feed URLs, refresh manually, delete feeds
- **Zero dependencies to run** — single binary + SQLite file
- **Server-rendered HTML** — works in any browser

---

## Getting Started

### Prerequisites

- [Go 1.21+](https://go.dev/dl/)
- GCC (required for SQLite via `go-sqlite3`)

### Install & Run

```bash
# Clone the repo
git clone https://github.com/tmarcum6/fed..git
cd fed./

# Install dependencies
go mod tidy

# Run
go run .
```

Then open [http://localhost:8080](http://localhost:8080) in your browser.

---

## Project Structure

```
fed/
├── main.go               # Entry point, router, handlers
├── go.mod
├── rss.db                # Auto-created SQLite database
├── db/
│   ├── db.go             # DB init & table creation
│   ├── feeds.go          # Feed queries
│   └── articles.go       # Article queries
├── handlers/
│   ├── discover.go        # RSS feed discovery on websites
│   └── import.go         # CSV feed import
├── models/
│   └── feed.go           # Feed & Article types
├── poller/
│   └── poller.go         # Background feed polling
└── templates/
    ├── templates.go       # Template loader & renderer
    ├── layout.html       # Base layout with sidebar
    ├── index.html        # Homepage / all articles
    ├── feed.html         # Single feed view
    └── feeds.html        # Feed management page
```

---

## Routes

| Method | Path | Description |
|--------|------|-------------|
| `GET` | `/` | Homepage — all articles |
| `GET` | `/?unread=true` | Unread articles only |
| `GET` | `/?hidden=true` | Hidden articles |
| `POST` | `/feeds` | Add a new feed |
| `GET` | `/feeds` | Manage feeds page |
| `POST` | `/feeds/import` | Import feeds from CSV |
| `GET` | `/feeds/discover` | Discover feeds on a website |
| `GET` | `/feeds/{id}` | View articles for a feed |
| `GET` | `/feeds/{id}?unread=true` | View unread articles for a feed |
| `POST` | `/feeds/{id}/refresh` | Manually refresh a feed |
| `PUT` | `/feeds/{id}` | Update feed URL |
| `DELETE` | `/feeds/{id}` | Delete a feed |
| `POST` | `/articles/{id}/read` | Mark article as read |
| `POST` | `/articles/{id}/unread` | Mark article as unread |
| `POST` | `/articles/{id}/hidden` | Hide an article |
| `POST` | `/articles/{id}/unhide` | Unhide an article |

---

## Testing

```bash
# Run all tests
go test ./...

# Run with verbose output
go test ./... -v

# Check coverage
go test ./... -cover
```

Tests use SQLite's `:memory:` mode — no files created, no cleanup needed.

---

## Tech Stack

| Layer | Tech |
|-------|------|
| Language | [Go](https://go.dev) |
| Router | [chi](https://github.com/go-chi/chi) |
| RSS Parsing | [gofeed](https://github.com/mmcdole/gofeed) |
| Database | [SQLite](https://sqlite.org) via [go-sqlite3](https://github.com/mattn/go-sqlite3) |
| Templates | [html/template](https://pkg.go.dev/html/template) |
| CSS | [Tailwind CSS](https://tailwindcss.com) |

---

## Roadmap

- [ ] OPML import / export
- [ ] Feed categories / folders
- [ ] Full-text search
- [ ] Article preview pane
- [ ] Dark / light theme toggle
- [ ] Docker support
- [ ] Tracking for Read and Saved
- [ ] Social
- [ ] Custom feeds
- [ ] Hot-reload

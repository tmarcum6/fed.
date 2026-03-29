# 📡 gRoSS

> **Go RSS** — a self-hosted feed reader that doesn't get in your way.

```
 ██████╗ ██████╗  ██████╗ ███████╗███████╗
██╔════╝ ██╔══██╗██╔═══██╗██╔════╝██╔════╝
██║  ███╗██████╔╝██║   ██║███████╗███████╗
██║   ██║██╔══██╗██║   ██║╚════██║╚════██║
╚██████╔╝██║  ██║╚██████╔╝███████║███████║
 ╚═════╝ ╚═╝  ╚═╝ ╚═════╝ ╚══════╝╚══════╝
```

**gRoSS** is a lightweight, self-hosted RSS feed reader built in Go. No Electron. No Node. No nonsense. Just a single binary, a SQLite database, and your feeds.

---

## ✨ Features

- 📥 **Add any RSS or Atom feed** via URL
- 🔄 **Auto-polling** — feeds refresh every 15 minutes in the background
- ✅ **Mark articles as read** — per article, right from the feed
- 🗂️ **Per-feed views** — browse articles from a single source
- 🔍 **Unread filter** — focus on what you haven't read yet
- 🪶 **Zero dependencies to run** — single binary + SQLite file
- 🖥️ **Server-rendered HTML** — works in any browser, no JS framework

---

## 🚀 Getting Started

### Prerequisites

- [Go 1.21+](https://go.dev/dl/)
- GCC (required for SQLite via `go-sqlite3`) — or swap to [`modernc.org/sqlite`](https://pkg.go.dev/modernc.org/sqlite) for a pure Go alternative

### Install & Run

```bash
# Clone the repo
git clone https://github.com/tmarcum6/gross.git
cd gross

# Install dependencies
go mod tidy

# Run
go run .
```

Then open [http://localhost:8080](http://localhost:8080) in your browser.

---

## 🗂️ Project Structure

```
go-rss/
├── main.go               # Entry point, router, handlers
├── go.mod
├── rss.db                # Auto-created SQLite database
├── db/
│   ├── db.go             # DB init & table creation
│   ├── feeds.go          # Feed queries
│   └── articles.go       # Article queries
├── models/
│   └── feed.go           # Feed & Article types
├── poller/
│   └── poller.go         # Background feed polling
└── templates/
    ├── templates.go      # Template loader & renderer
    ├── layout.html       # Base layout with sidebar
    ├── index.html        # Homepage / all articles
    └── feed.html         # Single feed view
```

---

## 🛣️ Routes

| Method | Path | Description |
|--------|------|-------------|
| `GET` | `/` | Homepage — all articles |
| `GET` | `/?unread=true` | Unread articles only |
| `POST` | `/feeds` | Add a new feed |
| `GET` | `/feeds/{id}` | View articles for a feed |
| `POST` | `/articles/{id}/read` | Mark an article as read |

---

## 🧪 Testing

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

## 🧰 Tech Stack

| Layer | Tech |
|-------|------|
| Language | [Go](https://go.dev) |
| Router | [chi](https://github.com/go-chi/chi) |
| RSS Parsing | [gofeed](https://github.com/mmcdole/gofeed) |
| Database | [SQLite](https://sqlite.org) via [go-sqlite3](https://github.com/mattn/go-sqlite3) |
| Templates | Go `html/template` |

---

## 🔮 Roadmap

- [ ] OPML import / export
- [ ] Feed categories / folders
- [ ] Full-text search
- [ ] Article preview pane
- [ ] Dark / light theme toggle
- [ ] Docker support
- [ ] Tailwind CSS

---

## 🤝 Contributing

PRs welcome.

```bash
git checkout -b feature/your-feature
git commit -m "add your feature"
git push origin feature/your-feature
# open a PR
```

---

## 📄 License

MIT — do whatever you want with it.

---

<p align="center">Built with Go & caffeine ☕</p>

package db

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

var DB *sql.DB

func Init(path string) {
	var err error
	DB, err = sql.Open("sqlite3", path)
	if err != nil {
		log.Fatal("failed to open db:", err)
	}

	if err = DB.Ping(); err != nil {
		log.Fatal("failed to connect to db:", err)
	}

	createTables()
	log.Println("database ready")
}

func createTables() {
	_, err := DB.Exec(`
		CREATE TABLE IF NOT EXISTS feeds (
		    id          INTEGER PRIMARY KEY AUTOINCREMENT,
		    url         TEXT UNIQUE NOT NULL,
		    title       TEXT,
		    last_fetched DATETIME
		);

		CREATE TABLE IF NOT EXISTS articles (
		    id          INTEGER PRIMARY KEY AUTOINCREMENT,
		    feed_id     INTEGER NOT NULL,
		    title       TEXT,
		    link        TEXT UNIQUE NOT NULL,
		    description TEXT,
		    published   DATETIME,
		    read        BOOLEAN DEFAULT 0,
		    hidden      BOOLEAN DEFAULT 0,
		    FOREIGN KEY (feed_id) REFERENCES feeds(id)
		);
	`)
	if err != nil {
		log.Fatal("failed to create tables:", err)
	}
}

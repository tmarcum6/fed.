package db

import (
	"gross/models"
	"time"
)

func InsertFeed(url, title string) (int64, error) {
	res, err := DB.Exec(
		`INSERT OR IGNORE INTO feeds (url, title, last_fetched) VALUES (?, ?, ?)`,
		url, title, time.Now(),
	)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func DeleteFeed(url string) (int64, error) {
	res, err := DB.Exec(
		`DELETE FROM feeds WHERE url = (?)`,
		url,
	)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

func GetAllFeeds() ([]models.Feed, error) {
	rows, err := DB.Query(`SELECT id, url, title, last_fetched FROM feeds`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var feeds []models.Feed
	for rows.Next() {
		var f models.Feed
		rows.Scan(&f.ID, &f.URL, &f.Title, &f.LastFetched)
		feeds = append(feeds, f)
	}
	return feeds, nil
}

func UpdateLastFetched(feedID int) error {
	_, err := DB.Exec(
		`UPDATE feeds SET last_fetched = ? WHERE id = ?`,
		time.Now(), feedID,
	)
	return err
}

func DeleteFeedByID(id int) error {
	_, err := DB.Exec(`DELETE FROM feeds WHERE id = ?`, id)
	return err
}

func UpdateFeedURL(id int, url, title string) error {
	_, err := DB.Exec(
		`UPDATE feeds SET url = ?, title = ? WHERE id = ?`,
		url, title, id,
	)
	return err
}

func GetFeedStats(feedID int) (total int, unread int, err error) {
	var e error
	row := DB.QueryRow(`SELECT COUNT(*) FROM articles WHERE feed_id = ?`, feedID)
	if e = row.Scan(&total); e != nil {
		return 0, 0, e
	}
	row = DB.QueryRow(`SELECT COUNT(*) FROM articles WHERE feed_id = ? AND read = 0 AND hidden = 0`, feedID)
	if e = row.Scan(&unread); e != nil {
		return 0, 0, e
	}
	return total, unread, nil
}

func GetFeedByID(id int) (*models.Feed, error) {
	var f models.Feed
	err := DB.QueryRow(`SELECT id, url, title, last_fetched FROM feeds WHERE id = ?`, id).
		Scan(&f.ID, &f.URL, &f.Title, &f.LastFetched)
	if err != nil {
		return nil, err
	}
	return &f, nil
}

func GetAllFeedURLs() ([]string, error) {
	rows, err := DB.Query(`SELECT url FROM feeds`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var urls []string
	for rows.Next() {
		var url string
		if err := rows.Scan(&url); err != nil {
			return nil, err
		}
		urls = append(urls, url)
	}
	return urls, nil
}

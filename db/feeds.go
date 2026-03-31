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

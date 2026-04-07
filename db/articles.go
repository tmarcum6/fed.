package db

import (
	"gross/models"
)

func InsertArticle(a models.Article) error {
	_, err := DB.Exec(`
        INSERT OR IGNORE INTO articles (feed_id, title, link, description, published)
        VALUES (?, ?, ?, ?, ?)`,
		a.FeedID, a.Title, a.Link, a.Description, a.Published,
	)
	return err
}

func GetArticlesByFeed(feedID int) ([]models.Article, error) {
	rows, err := DB.Query(
		`SELECT id, feed_id, title, link, description, published, read
	         FROM articles WHERE feed_id = ? ORDER BY published DESC`,
		feedID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var articles []models.Article
	for rows.Next() {
		var a models.Article
		rows.Scan(&a.ID, &a.FeedID, &a.Title, &a.Link, &a.Description, &a.Published, &a.Read)
		articles = append(articles, a)
	}
	return articles, nil
}

func MarkAsRead(articleID int) error {
	_, err := DB.Exec(`UPDATE articles SET read = 1 WHERE id = ?`, articleID)
	return err
}

func MarkAsHidden(articleID int) error {
	_, err := DB.Exec(`UPDATE articles SET hidden = 1 WHERE id = ?`, articleID)
	return err
}

func UnhideArticle(articleID int) error {
	_, err := DB.Exec(`UPDATE articles SET hidden = 0 WHERE id = ?`, articleID)
	return err
}

func MarkAsUnread(articleID int) error {
	_, err := DB.Exec(`UPDATE articles SET read = 0 WHERE id = ?`, articleID)
	return err
}

func DeleteArticlesByFeed(feedID int) error {
	_, err := DB.Exec(`DELETE FROM articles WHERE feed_id = ?`, feedID)
	return err
}

func GetArticles(feedID string, unreadOnly bool, hiddenOnly bool) ([]models.Article, error) {
	query := `SELECT id, feed_id, title, link, description, published, read, hidden
              FROM articles WHERE 1=1`
	args := []any{}

	if feedID != "" {
		query += " AND feed_id = ?"
		args = append(args, feedID)
	}
	if unreadOnly {
		query += " AND read = 0 AND hidden = 0"
	} else if hiddenOnly {
		query += " AND hidden = 1"
	} else {
		query += " AND hidden = 0"
	}

	query += " ORDER BY published DESC LIMIT 100"

	rows, err := DB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var articles []models.Article
	for rows.Next() {
		var a models.Article
		rows.Scan(&a.ID, &a.FeedID, &a.Title, &a.Link, &a.Description, &a.Published, &a.Read, &a.Hidden)
		articles = append(articles, a)
	}
	return articles, nil
}

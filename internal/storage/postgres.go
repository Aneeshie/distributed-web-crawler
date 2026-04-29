package storage

import (
	"database/sql"
	"log"
)

//THE ONLY RESPONSIBILITY IS CONNECT TO POSTGRESQL

func NewPostgresDB(connStr string) *sql.DB {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("failed to open DB:", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal("failed to connect DB:", err)
	}

	return db
}

func SavePage(db *sql.DB, url, title string, status int) error {
	query := `
	INSERT INTO pages (url, title, status_code)
	VALUES ($1, $2, $3)
	ON CONFLICT (url) DO UPDATE
	SET title = EXCLUDED.title,
	    status_code = EXCLUDED.status_code,
	    crawled_at = NOW()
	`

	_, err := db.Exec(query, url, title, status)
	return err
}

func SaveLinks(db *sql.DB, source string, links []string) error {
	query := `
	INSERT INTO links (source_url, target_url)
	VALUES ($1, $2)
	ON CONFLICT DO NOTHING
	`

	for _, link := range links {
		_, err := db.Exec(query, source, link)
		if err != nil {
			return err
		}
	}

	return nil
}


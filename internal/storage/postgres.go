package storage

import (
	"database/sql"
	"log"
	"net/url"

	_ "github.com/lib/pq"
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

func SavePage(db *sql.DB, rawURL, title string, status int) error {

	parsed, err := url.Parse(rawURL)
	if err != nil {
		return err
	}

	domain := parsed.Host

	query := `
	INSERT INTO pages (url, domain, title, status_code)
	VALUES ($1, $2, $3, $4)
	ON CONFLICT (url) DO UPDATE
	SET domain = EXCLUDED.domain,
	    title = EXCLUDED.title,
	    status_code = EXCLUDED.status_code,
	    crawled_at = NOW()
	`

	_, err = db.Exec(query, rawURL, domain, title, status)
	return err
}

// PageExists reports whether rawURL already has a row in pages (already crawled at least once).
func PageExists(db *sql.DB, rawURL string) (bool, error) {
	const q = `SELECT EXISTS (SELECT 1 FROM pages WHERE url = $1 LIMIT 1)`
	var exists bool
	err := db.QueryRow(q, rawURL).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
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

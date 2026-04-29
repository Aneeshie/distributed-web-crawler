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


package storage

import (
	"database/sql"
	"log"
)

func CreateTables(db *sql.DB) {
	pagesTable := `
	CREATE TABLE IF NOT EXISTS pages (
		id SERIAL PRIMARY KEY,
		url TEXT UNIQUE NOT NULL,
		title TEXT,
		status_code INT,
		crawled_at TIMESTAMP DEFAULT NOW()
	);`

	linksTable := `
	CREATE TABLE IF NOT EXISTS links (
		id SERIAL PRIMARY KEY,
		source_url TEXT NOT NULL,
		target_url TEXT NOT NULL
	);`

	_, err := db.Exec(pagesTable)
	if err != nil {
		log.Fatal("failed creating pages table:", err)
	}

	_, err = db.Exec(linksTable)
	if err != nil {
		log.Fatal("failed creating links table:", err)
	}

	log.Println("tables ready")
}

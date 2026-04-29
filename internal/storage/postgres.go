package storage

import "database/sql"

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

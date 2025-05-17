package database

import (
	"database/sql"
	"log"
	"os"

	_ "modernc.org/sqlite"
)

// DB is the database connection
var DB *sql.DB

// InitDB initializes the SQLite database
func InitDB() error {
	var err error
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "job_tracker.db"
	}
	log.Printf("Opening database at %s", dbPath)
	DB, err = sql.Open("sqlite", dbPath)
	if err != nil {
		return err
	}
	DB.SetMaxOpenConns(10)
	DB.SetMaxIdleConns(5)

	// Create table if it doesn't exist
	createTableQuery := `
		CREATE TABLE IF NOT EXISTS job_applications (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			company TEXT,
			position TEXT,
			application_date TEXT,
			status TEXT
		);
	`
	_, err = DB.Exec(createTableQuery)
	if err != nil {
		return err
	}

	return nil
}

// CloseDB closes the database connection
func CloseDB() {
	if DB != nil {
		if err := DB.Close(); err != nil {
			log.Printf("Error closing database: %v", err)
		}
	}
}

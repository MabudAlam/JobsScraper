package store

import (
	"database/sql"
	"os"

	_ "modernc.org/sqlite"
)

var db *sql.DB

func InitDB() error {
	dbPath := os.Getenv("DATABASE_PATH")
	if dbPath == "" {
		dbPath = "./jobs.db"
	}

	var err error
	db, err = sql.Open("sqlite", dbPath)
	if err != nil {
		return err
	}

	if err = db.Ping(); err != nil {
		return err
	}

	if err = createTables(); err != nil {
		return err
	}

	return nil
}

func createTables() error {
	schema := `
	CREATE TABLE IF NOT EXISTS companies (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		ashby_slug TEXT NOT NULL UNIQUE,
		last_scraped_at TIMESTAMP,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS jobs (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		job_id TEXT NOT NULL,
		company TEXT NOT NULL,
		title TEXT NOT NULL,
		location TEXT,
		team TEXT,
		department TEXT,
		employment_type TEXT,
		remote INTEGER NOT NULL DEFAULT 0,
		description TEXT,
		apply_url TEXT,
		job_url TEXT,
		published_at TIMESTAMP,
		scraped_at TIMESTAMP NOT NULL,
		compensation_summary TEXT,
		content_hash TEXT NOT NULL,
		is_active INTEGER NOT NULL DEFAULT 1,
		status TEXT DEFAULT 'new',
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		UNIQUE(company, job_id)
	);

	CREATE INDEX IF NOT EXISTS idx_jobs_active ON jobs(is_active);
	CREATE INDEX IF NOT EXISTS idx_jobs_company ON jobs(company);

	CREATE TABLE IF NOT EXISTS job_snapshots (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		job_id TEXT NOT NULL,
		company TEXT NOT NULL,
		content_hash TEXT NOT NULL,
		snapshot_data TEXT NOT NULL,
		captured_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);

	CREATE INDEX IF NOT EXISTS idx_snapshots_job ON job_snapshots(company, job_id);
	`

	_, err := db.Exec(schema)
	return err
}

func GetDB() *sql.DB {
	return db
}

func CloseDB() error {
	if db != nil {
		return db.Close()
	}
	return nil
}

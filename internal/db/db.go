package db

import (
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite"
)

// Open opens (or creates) the SQLite database at dataDir/stories.db
// and runs the schema migration.
func Open(dataDir string) (*sql.DB, error) {
	dsn := fmt.Sprintf("file:%s/stories.db?_journal_mode=WAL&_foreign_keys=on", dataDir)
	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, fmt.Errorf("open db: %w", err)
	}

	db.SetMaxOpenConns(1) // SQLite is single-writer; keep it simple.

	if err := migrate(db); err != nil {
		db.Close()
		return nil, fmt.Errorf("migrate: %w", err)
	}
	return db, nil
}

func migrate(db *sql.DB) error {
	const query = `
	CREATE TABLE IF NOT EXISTS stories (
		id           INTEGER  PRIMARY KEY AUTOINCREMENT,
		title        TEXT     NOT NULL,
		cover_image  TEXT     NOT NULL DEFAULT '',
		author       TEXT     NOT NULL,
		content      TEXT     NOT NULL DEFAULT '',
		ai_generated INTEGER  NOT NULL DEFAULT 0,
		size         TEXT     NOT NULL DEFAULT 'small',
		views        INTEGER  NOT NULL DEFAULT 0,
		created_at   DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at   DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
	);`
	_, err := db.Exec(query)
	return err
}

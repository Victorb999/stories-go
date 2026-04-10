package db

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

// Open connects to the PostgreSQL database using the DATABASE_URL
// and runs the schema migration.
func Open(databaseURL string) (*sql.DB, error) {
	if databaseURL == "" {
		return nil, fmt.Errorf("DATABASE_URL is not set")
	}

	db, err := sql.Open("pgx", databaseURL)
	if err != nil {
		return nil, fmt.Errorf("open db: %w", err)
	}

	// Ping to ensure connection is valid
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		db.Close()
		return nil, fmt.Errorf("ping db: %w", err)
	}

	if err := migrate(db); err != nil {
		db.Close()
		return nil, fmt.Errorf("migrate: %w", err)
	}
	return db, nil
}

func migrate(db *sql.DB) error {
	const query = `
	CREATE TABLE IF NOT EXISTS stories (
		id           BIGSERIAL PRIMARY KEY,
		title        TEXT     NOT NULL,
		cover_image  TEXT     NOT NULL DEFAULT '',
		author       TEXT     NOT NULL,
		content      TEXT     NOT NULL DEFAULT '',
		ai_generated BOOLEAN  NOT NULL DEFAULT FALSE,
		size         TEXT     NOT NULL DEFAULT 'small',
		views        BIGINT   NOT NULL DEFAULT 0,
		created_at   TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at   TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
	);`
	_, err := db.Exec(query)
	return err
}

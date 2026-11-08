package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

// open is the shared logic for every database: open → pool settings → ping.
// Lowercase = unexported: only code inside package database can call it.
// Fix a bug here once, and MySQL AND Postgres both get the fix.
func open(driver, dsn string) (*sql.DB, error) {
	db, err := sql.Open(driver, dsn)
	if err != nil {
		return nil, fmt.Errorf("%s open: %w", driver, err)
	}

	// Pool settings
	db.SetMaxOpenConns(25)                 // max connections in use at once
	db.SetMaxIdleConns(25)                 // connections kept ready for reuse
	db.SetConnMaxLifetime(5 * time.Minute) // recycle connections after 5 min

	// Actually connect, and fail fast if the database is unreachable
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		db.Close()
		return nil, fmt.Errorf("%s ping: %w", driver, err)
	}

	return db, nil
}

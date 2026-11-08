package database

import (
	"database/sql"

	_ "github.com/jackc/pgx/v5/stdlib" // registers the "pgx" driver with database/sql
)

// NewPostgresConnection opens a Postgres connection pool and checks that it works.
// dsn format: postgres://user:password@host:5432/dbname?sslmode=disable
func NewPostgresConnection(dsn string) (*sql.DB, error) {
	return open("pgx", dsn)
}

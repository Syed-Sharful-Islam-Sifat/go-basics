package user

import (
	"context"
	"database/sql"
)

// PostgresRepo stores users in Postgres.
// Same shape as MySQLRepo — only the SQL inside the methods differs.
type PostgresRepo struct {
	db *sql.DB
}

// NewPostgresRepo creates a Postgres-backed user repository.
func NewPostgresRepo(db *sql.DB) *PostgresRepo {
	return &PostgresRepo{db: db}
}

func (r *PostgresRepo) CreateTable(ctx context.Context) error {
	_, err := r.db.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS users (
			id    SERIAL PRIMARY KEY,               -- Postgres: SERIAL (MySQL: AUTO_INCREMENT)
			name  VARCHAR(100) NOT NULL,
			email VARCHAR(255) NOT NULL
		)`)
	return err
}

func (r *PostgresRepo) Create(ctx context.Context, name, email string) (int, error) {
	// Postgres: $1, $2 placeholders (MySQL: ?)
	// Postgres has no LastInsertId(), so ask for the id back with RETURNING.
	var id int
	err := r.db.QueryRowContext(ctx,
		"INSERT INTO users (name, email) VALUES ($1, $2) RETURNING id", name, email).
		Scan(&id)
	return id, err
}

func (r *PostgresRepo) GetByID(ctx context.Context, id int) (User, error) {
	var u User
	err := r.db.QueryRowContext(ctx,
		"SELECT id, name, email FROM users WHERE id = $1", id).
		Scan(&u.ID, &u.Name, &u.Email)
	return u, err
}

// Compile-time check: fails to build if PostgresRepo is missing any Repository method.
var _ Repository = (*PostgresRepo)(nil)

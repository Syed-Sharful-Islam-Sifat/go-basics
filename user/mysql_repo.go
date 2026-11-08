package user

import (
	"context"
	"database/sql"
)

type MySQLRepo struct {
	db *sql.DB
}
// db *sql.DB means "share the one connection pool"

// remember that *sql.DB is the taxi stand with 25 taxis
/*
  -- there should be one taxi stand for the whole app
  -- main creates it once with NewMySQLConnection
  -- The repo , the service and future repos all need to use that same stand

*/
func NewMySQLRepo(db *sql.DB) *MySQLRepo{
	return &MySQLRepo{db:db}
}

// &MySQLRepo{db:db} this means sharing the one 

func (r *MySQLRepo) CreateTable(ctx context.Context) error {
	_, err := r.db.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS users (
			id    INT AUTO_INCREMENT PRIMARY KEY,
			name  VARCHAR(100) NOT NULL,
			email VARCHAR(255) NOT NULL
		)`)
	return err
}

func (r *MySQLRepo) Create(ctx context.Context, name, email string) (int, error) {
	res, err := r.db.ExecContext(ctx,
		"INSERT INTO users (name, email) VALUES (?, ?)", name, email)
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	return int(id), err
}

func (r *MySQLRepo) GetByID(ctx context.Context, id int) (User, error) {
	var u User
	err := r.db.QueryRowContext(ctx,
		"SELECT id, name, email FROM users WHERE id = ?", id).
		Scan(&u.ID, &u.Name, &u.Email)
	return u, err
}

// Compile-time check: fails to build if MySQLRepo is missing any Repository method.
var _ Repository = (*MySQLRepo)(nil)

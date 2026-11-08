package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"

	"pggs/dbapp/database"
	"pggs/dbapp/user"
)

// getEnv reads an env variable, or returns a fallback if it's empty.
func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func main() {
	ctx := context.Background()

	// ─────────────────────────────────────────────────────────────
	// The ONLY place in the whole app that knows which database is used.
	// Switch with:  DB_DRIVER=postgres go run ./dbapp
	// ─────────────────────────────────────────────────────────────
	var (
		db   *sql.DB
		repo user.Repository // ← the interface: can hold MySQLRepo OR PostgresRepo
		err  error
	)

	driver := getEnv("DB_DRIVER", "mysql")
	switch driver {
	case "mysql":
		db, err = database.NewMySQLConnection(
			getEnv("MYSQL_DSN", "root:secret@tcp(localhost:3306)/shop?parseTime=true"))
		if err == nil {
			repo = user.NewMySQLRepo(db)
		}
	case "postgres":
		db, err = database.NewPostgresConnection(
			getEnv("POSTGRES_DSN", "postgres://postgres:secret@localhost:5432/shop?sslmode=disable"))
		if err == nil {
			repo = user.NewPostgresRepo(db)
		}
	default:
		log.Fatalf("unknown DB_DRIVER %q (use mysql or postgres)", driver)
	}
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	fmt.Println("✅ connected using", driver)

	// ─────────────────────────────────────────────────────────────
	// Everything below is written ONCE and works for every database.
	// No "MySQL" or "Postgres" anywhere — just the service.
	// ─────────────────────────────────────────────────────────────
	svc := user.NewService(repo)

	if err := svc.Setup(ctx); err != nil {
		log.Fatal(err)
	}

	id, err := svc.Register(ctx, "tm", "tm@example.com")
	if err != nil {
		log.Fatal(err)
	}

	u, err := svc.Get(ctx, id)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("saved and loaded: %+v\n", u)

	// Business rules apply no matter which database is behind the service
	if _, err := svc.Register(ctx, "bad", "no-at-sign"); err != nil {
		fmt.Println("rejected as expected:", err)
	}
}

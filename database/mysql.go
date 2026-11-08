package database // it's must so other package can use it.

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

// database/sql - Go's standard database API. it gives you sql.Open and the *sql DB type. it works with any sql database but does not know how to talk
// context lets you set a timeout on the connection check , so a dead server does not hang you up
// Duration like 5*time.Second and 5*time.Minute ,  for the timeout and the pool settings
// fmt.Errorf to wrap errors with a message "mysql ping: ..."
// "github.com/go-sql-driver/mysql" The MySQL driver , the piece that actually knows how to talk to MySQL

// (context, time and fmt moved to open.go — the shared open → pool → ping logic)

// NewMySQLConnection opens a MySQL connection pool and checks that it works.
// dsn format: user:password@tcp(host:3306)/dbname?parseTime=true
func NewMySQLConnection(dsn string) (*sql.DB, error) {
	return open("mysql", dsn)
}

// New... A Go convention: functions that create and return something start with New (like a constructor)

// Capital N Exported , so main can call database.NewMySQLConnection()

// dsn string this is a connection string that has passed

// *sql.DB The connection pool that the rest of your app will use

// error: Go has no exceptions. Functions that can fail return an error as the last value


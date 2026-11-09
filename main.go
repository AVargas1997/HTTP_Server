package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"

	// Registers the "pgx" PostgreSQL driver with database/sql.
	_ "github.com/jackc/pgx/v5/stdlib"
)

// main connects to PostgreSQL, makes sure the tasks table exists,
// and starts the HTTP server.
func main() {
	// The connection string comes from an environment variable so the
	// password is never written in the code.
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL is not set, for example: " +
			"postgres://postgres:password@localhost:5432/tasks?sslmode=disable")
	}

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// sql.Open doesn't actually connect, so ping to fail fast if the
	// database is down or the connection string is wrong.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		log.Fatal("could not connect to database: ", err)
	}

	store := NewPostgresStore(db)
	if err := store.Migrate(); err != nil {
		log.Fatal("could not create tasks table: ", err)
	}

	server := NewServer(store)

	log.Println("Server running at http://localhost:8080")

	err = http.ListenAndServe("localhost:8080", server.routes())
	if err != nil {
		log.Fatal(err)
	}
}

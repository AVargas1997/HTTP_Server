package main

import (
	"database/sql"
	"os"
	"testing"
)

// TestPostgresStore runs against a real database. It is skipped unless
// TEST_DATABASE_URL is set, so `go test` still works without Postgres.
// WARNING: it deletes everything in the tasks table, so point it at a
// separate test database, not the one you use with the server.
func TestPostgresStore(t *testing.T) {
	dsn := os.Getenv("TEST_DATABASE_URL")
	if dsn == "" {
		t.Skip("TEST_DATABASE_URL not set, skipping PostgreSQL test")
	}

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	store := NewPostgresStore(db)
	if err := store.Migrate(); err != nil {
		t.Fatal(err)
	}
	if _, err := db.Exec(`TRUNCATE tasks RESTART IDENTITY`); err != nil {
		t.Fatal(err)
	}

	created, err := store.Create("Learn Postgres")
	if err != nil {
		t.Fatal(err)
	}
	if created.ID != 1 || created.Title != "Learn Postgres" || created.Completed {
		t.Errorf("unexpected task: %+v", created)
	}

	got, err := store.Get(created.ID)
	if err != nil || got != created {
		t.Errorf("Get returned %+v, %v", got, err)
	}

	tasks, err := store.List()
	if err != nil || len(tasks) != 1 {
		t.Errorf("List returned %v, %v", tasks, err)
	}

	if err := store.Delete(created.ID); err != nil {
		t.Fatal(err)
	}
	if _, err := store.Get(created.ID); err != ErrNotFound {
		t.Errorf("expected ErrNotFound after delete, got %v", err)
	}
	if err := store.Delete(created.ID); err != ErrNotFound {
		t.Errorf("expected ErrNotFound deleting twice, got %v", err)
	}
}

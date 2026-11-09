package main

import "errors"

// ErrNotFound is returned by a TaskStore when no task has the requested ID.
var ErrNotFound = errors.New("task not found")

// TaskStore is anything that can save and load tasks.
// The handlers only talk to this interface, so the same handlers work
// with PostgreSQL (PostgresStore) or plain memory (MemoryStore, used in tests).
type TaskStore interface {
	Create(title string) (Task, error)
	List() ([]Task, error)
	Get(id int) (Task, error)
	Delete(id int) error
}

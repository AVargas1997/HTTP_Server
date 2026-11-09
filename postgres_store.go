package main

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

// PostgresStore saves tasks in a PostgreSQL table, so they survive restarts.
// *sql.DB is a connection pool that is safe to use from many requests at
// once, so no mutex is needed here.
type PostgresStore struct {
	db *sql.DB
}

// NewPostgresStore wraps an open database connection pool.
func NewPostgresStore(db *sql.DB) *PostgresStore {
	return &PostgresStore{db: db}
}

// queryTimeout stops a slow or stuck database from hanging a request forever.
const queryTimeout = 5 * time.Second

// Migrate creates the tasks table if it doesn't exist yet.
func (p *PostgresStore) Migrate() error {
	ctx, cancel := context.WithTimeout(context.Background(), queryTimeout)
	defer cancel()

	_, err := p.db.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS tasks (
			id        SERIAL PRIMARY KEY,
			title     TEXT    NOT NULL,
			completed BOOLEAN NOT NULL DEFAULT FALSE
		)`)
	return err
}

func (p *PostgresStore) Create(title string) (Task, error) {
	ctx, cancel := context.WithTimeout(context.Background(), queryTimeout)
	defer cancel()

	// $1 is a placeholder: the driver sends title separately from the SQL,
	// which prevents SQL injection.
	var task Task
	err := p.db.QueryRowContext(ctx,
		`INSERT INTO tasks (title) VALUES ($1)
		 RETURNING id, title, completed`,
		title,
	).Scan(&task.ID, &task.Title, &task.Completed)

	return task, err
}

func (p *PostgresStore) List() ([]Task, error) {
	ctx, cancel := context.WithTimeout(context.Background(), queryTimeout)
	defer cancel()

	rows, err := p.db.QueryContext(ctx,
		`SELECT id, title, completed FROM tasks ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tasks := []Task{}
	for rows.Next() {
		var task Task
		if err := rows.Scan(&task.ID, &task.Title, &task.Completed); err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}

	return tasks, rows.Err()
}

func (p *PostgresStore) Get(id int) (Task, error) {
	ctx, cancel := context.WithTimeout(context.Background(), queryTimeout)
	defer cancel()

	var task Task
	err := p.db.QueryRowContext(ctx,
		`SELECT id, title, completed FROM tasks WHERE id = $1`,
		id,
	).Scan(&task.ID, &task.Title, &task.Completed)

	if errors.Is(err, sql.ErrNoRows) {
		return Task{}, ErrNotFound
	}

	return task, err
}

func (p *PostgresStore) Delete(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), queryTimeout)
	defer cancel()

	result, err := p.db.ExecContext(ctx, `DELETE FROM tasks WHERE id = $1`, id)
	if err != nil {
		return err
	}

	deleted, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if deleted == 0 {
		return ErrNotFound
	}

	return nil
}

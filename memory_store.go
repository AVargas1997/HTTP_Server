package main

import (
	"sort"
	"sync"
)

// MemoryStore keeps tasks in a map. Data is lost when the program stops.
// It is used by the tests so they don't need a database.
type MemoryStore struct {
	mu     sync.Mutex
	tasks  map[int]Task
	nextID int
}

// NewMemoryStore creates an empty in-memory store.
func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		tasks:  make(map[int]Task),
		nextID: 1,
	}
}

func (m *MemoryStore) Create(title string) (Task, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	task := Task{ID: m.nextID, Title: title}
	m.tasks[task.ID] = task
	m.nextID++

	return task, nil
}

func (m *MemoryStore) List() ([]Task, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	tasks := make([]Task, 0, len(m.tasks))
	for _, task := range m.tasks {
		tasks = append(tasks, task)
	}

	// Sort by ID so the order matches the Postgres store.
	sort.Slice(tasks, func(i, j int) bool { return tasks[i].ID < tasks[j].ID })

	return tasks, nil
}

func (m *MemoryStore) Get(id int) (Task, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	task, exists := m.tasks[id]
	if !exists {
		return Task{}, ErrNotFound
	}

	return task, nil
}

func (m *MemoryStore) Delete(id int) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.tasks[id]; !exists {
		return ErrNotFound
	}

	delete(m.tasks, id)
	return nil
}

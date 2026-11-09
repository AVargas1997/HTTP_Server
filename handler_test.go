package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestCreateTask(t *testing.T) {
	store := NewMemoryStore()
	server := NewServer(store)

	request := httptest.NewRequest(
		http.MethodPost,
		"/tasks",
		strings.NewReader(`{"title":" Learn Go "}`),
	)
	response := httptest.NewRecorder()

	server.routes().ServeHTTP(response, request)

	if response.Code != http.StatusCreated {
		t.Fatalf("expected status 201, got %d", response.Code)
	}

	var task Task
	err := json.Unmarshal(response.Body.Bytes(), &task)
	if err != nil {
		t.Fatalf("could not decode response: %v", err)
	}

	if task.ID != 1 {
		t.Errorf("expected ID 1, got %d", task.ID)
	}

	if task.Title != "Learn Go" {
		t.Errorf("expected title %q, got %q", "Learn Go", task.Title)
	}

	if task.Completed {
		t.Error("expected a new task to be incomplete")
	}

	if response.Header().Get("Location") != "/tasks/1" {
		t.Errorf(
			"expected Location /tasks/1, got %q",
			response.Header().Get("Location"),
		)
	}

	if _, err := store.Get(task.ID); err != nil {
		t.Error("expected task to be saved")
	}
}

func TestCreateTaskEmptyTitle(t *testing.T) {
	store := NewMemoryStore()
	server := NewServer(store)

	request := httptest.NewRequest(
		http.MethodPost,
		"/tasks",
		strings.NewReader(`{"title":"   "}`),
	)
	response := httptest.NewRecorder()

	server.routes().ServeHTTP(response, request)

	if response.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", response.Code)
	}

	if tasks, _ := store.List(); len(tasks) != 0 {
		t.Error("expected no tasks to be saved")
	}
}

func TestListTasks(t *testing.T) {
	store := NewMemoryStore()
	server := NewServer(store)
	store.Create("Learn Go")

	request := httptest.NewRequest(http.MethodGet, "/tasks", nil)
	response := httptest.NewRecorder()

	server.routes().ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", response.Code)
	}

	var tasks []Task
	err := json.Unmarshal(response.Body.Bytes(), &tasks)
	if err != nil {
		t.Fatalf("could not decode response: %v", err)
	}

	if len(tasks) != 1 {
		t.Fatalf("expected 1 task, got %d", len(tasks))
	}

	if tasks[0].Title != "Learn Go" {
		t.Errorf("expected title %q, got %q", "Learn Go", tasks[0].Title)
	}
}

func TestGetTask(t *testing.T) {
	store := NewMemoryStore()
	server := NewServer(store)
	store.Create("Learn Go")

	request := httptest.NewRequest(http.MethodGet, "/tasks/1", nil)
	response := httptest.NewRecorder()

	server.routes().ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", response.Code)
	}

	var task Task
	err := json.Unmarshal(response.Body.Bytes(), &task)
	if err != nil {
		t.Fatalf("could not decode response: %v", err)
	}

	if task.ID != 1 || task.Title != "Learn Go" {
		t.Errorf("unexpected task: %+v", task)
	}
}

func TestGetTaskNotFound(t *testing.T) {
	store := NewMemoryStore()
	server := NewServer(store)

	request := httptest.NewRequest(http.MethodGet, "/tasks/99", nil)
	response := httptest.NewRecorder()

	server.routes().ServeHTTP(response, request)

	if response.Code != http.StatusNotFound {
		t.Fatalf("expected status 404, got %d", response.Code)
	}
}

func TestDeleteTask(t *testing.T) {
	store := NewMemoryStore()
	server := NewServer(store)
	store.Create("Learn Go")

	request := httptest.NewRequest(http.MethodDelete, "/tasks/1", nil)
	response := httptest.NewRecorder()

	server.routes().ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", response.Code)
	}

	if _, err := store.Get(1); err != ErrNotFound {
		t.Error("expected task to be deleted")
	}
}

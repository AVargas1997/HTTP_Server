package main

import (
	"net/http"
	"sync"
)

// Server holds the shared task storage.
type Server struct {
	mu     sync.Mutex
	tasks  map[int]Task
	nextID int
}

// NewServer creates a server with empty task storage.
func NewServer() *Server {
	return &Server{
		tasks:  make(map[int]Task),
		nextID: 1,
	}
}

// routes connects each URL to its handler.
func (s *Server) routes() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /tasks", s.createTask)
	mux.HandleFunc("GET /tasks", s.listTasks)
	mux.HandleFunc("GET /tasks/{id}", s.getTask)
	mux.HandleFunc("DELETE /tasks/{id}", s.deleteTask)

	return mux
}

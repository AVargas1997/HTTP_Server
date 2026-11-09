package main

import "net/http"

// Server holds the task store the handlers read from and write to.
type Server struct {
	store TaskStore
}

// NewServer creates a server that uses the given store.
func NewServer(store TaskStore) *Server {
	return &Server{store: store}
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

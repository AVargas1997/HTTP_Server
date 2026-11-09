package main

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"strings"
)

// createTask handles POST /tasks.
func (s *Server) createTask(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Title string `json:"title"`
	}

	r.Body = http.MaxBytesReader(w, r.Body, 1_048_576)

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	err := decoder.Decode(&input)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{
			"error": "Body must be a JSON object with a string title",
		})
		return
	}

	// Require exactly one JSON value in the request body.
	err = decoder.Decode(&struct{}{})
	if err != io.EOF {
		writeJSON(w, http.StatusBadRequest, map[string]string{
			"error": "Body must contain exactly one JSON object",
		})
		return
	}

	input.Title = strings.TrimSpace(input.Title)

	if input.Title == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{
			"error": "Title is required",
		})
		return
	}

	s.mu.Lock()

	task := Task{
		ID:        s.nextID,
		Title:     input.Title,
		Completed: false,
	}

	s.tasks[task.ID] = task
	s.nextID++

	s.mu.Unlock()

	w.Header().Set("Location", "/tasks/"+strconv.Itoa(task.ID))
	writeJSON(w, http.StatusCreated, task)
}

// listTasks handles GET /tasks.
func (s *Server) listTasks(w http.ResponseWriter, r *http.Request) {
	s.mu.Lock()

	tasks := make([]Task, 0, len(s.tasks))

	for _, task := range s.tasks {
		tasks = append(tasks, task)
	}

	s.mu.Unlock()

	writeJSON(w, http.StatusOK, tasks)
}

// getTask handles GET /tasks/{id}.
func (s *Server) getTask(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id < 1 {
		writeJSON(w, http.StatusBadRequest, map[string]string{
			"error": "Task ID must be a positive integer",
		})
		return
	}

	s.mu.Lock()
	task, exists := s.tasks[id]
	s.mu.Unlock()

	if !exists {
		writeJSON(w, http.StatusNotFound, map[string]string{
			"error": "Task not found",
		})
		return
	}

	writeJSON(w, http.StatusOK, task)
}

// deleteTask handles DELETE /tasks/{id}.
func (s *Server) deleteTask(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id < 1 {
		writeJSON(w, http.StatusBadRequest, map[string]string{
			"error": "Task ID must be a positive integer",
		})
		return
	}

	s.mu.Lock()

	_, exists := s.tasks[id]
	if exists {
		delete(s.tasks, id)
	}

	s.mu.Unlock()

	if !exists {
		writeJSON(w, http.StatusNotFound, map[string]string{
			"error": "Task not found",
		})
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{
		"message": "Task deleted",
	})
}

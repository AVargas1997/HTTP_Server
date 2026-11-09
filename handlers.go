package main

import (
	"encoding/json"
	"errors"
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

	task, err := s.store.Create(input.Title)
	if err != nil {
		serverError(w, err)
		return
	}

	w.Header().Set("Location", "/tasks/"+strconv.Itoa(task.ID))
	writeJSON(w, http.StatusCreated, task)
}

// listTasks handles GET /tasks.
func (s *Server) listTasks(w http.ResponseWriter, r *http.Request) {
	tasks, err := s.store.List()
	if err != nil {
		serverError(w, err)
		return
	}

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

	task, err := s.store.Get(id)
	if errors.Is(err, ErrNotFound) {
		writeJSON(w, http.StatusNotFound, map[string]string{
			"error": "Task not found",
		})
		return
	}
	if err != nil {
		serverError(w, err)
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

	err = s.store.Delete(id)
	if errors.Is(err, ErrNotFound) {
		writeJSON(w, http.StatusNotFound, map[string]string{
			"error": "Task not found",
		})
		return
	}
	if err != nil {
		serverError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{
		"message": "Task deleted",
	})
}

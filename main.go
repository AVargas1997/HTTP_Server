package main

import (
	"log"
	"net/http"
)

// / setup / starting point of the program... creates task storage, connects URLs to handles and listens for requests.
func main() {
	///creats the server and its storage
	server := &Server{
		/// empty task map to start
		tasks: make(map[int]Task),
		// first created task ID
		nextID: 1,
	}
	///creates rewuest router
	mux := http.NewServeMux()

	// assigns route to each handler
	mux.HandleFunc("POST /tasks", server.createTask)
	mux.HandleFunc("GET /tasks", server.listTasks)
	mux.HandleFunc("GET /tasks/{id}", server.getTask)
	mux.HandleFunc("DELETE /tasks/{id}", server.deleteTask)

	//announce the address that is about to be used
	log.Println("Server running at http://localhost:8080")

	// listens for requests then uses mux to route them
	//stays running while the server is serving
	err := http.ListenAndServe("localhost:8080", mux)
	if err != nil { //  if error is returnes , report and exit.
		log.Fatal(err)
	}
}

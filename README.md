# HTTP Server (Task API)

This is a small HTTP server I built in Go to learn how web servers and REST APIs work. It lets you create, view, and delete tasks (like a really simple to-do list) by sending HTTP requests.

The tasks are only saved in memory, so they go away when the server stops. I didn't add a database yet.

## What it does

| Method | URL | What it does |
|--------|-----|--------------|
| `POST` | `/tasks` | Creates a new task |
| `GET` | `/tasks` | Lists all tasks |
| `GET` | `/tasks/{id}` | Gets one task by its ID |
| `DELETE` | `/tasks/{id}` | Deletes a task by its ID |

A task looks like this:

```json
{
  "id": 1,
  "title": "Learn Go",
  "completed": false
}
```

## How to run it

You need Go 1.22 or newer (the routes use the `"GET /tasks/{id}"` style patterns, which were added in 1.22).

1. Clone the repo:

```bash
git clone https://github.com/AVargas1997/HTTP_Server.git
```

2. Go into the folder:

```bash
cd HTTP_Server
```

3. Start the server:

```bash
go run .
```

You should see `Server running at http://localhost:8080`.

## Trying it out

I used `curl` to test it. Open a second terminal while the server is running.

Create a task:

```bash
curl -X POST http://localhost:8080/tasks -d '{"title":"Learn Go"}'
```

List all tasks:

```bash
curl http://localhost:8080/tasks
```

Get one task:

```bash
curl http://localhost:8080/tasks/1
```

Delete a task:

```bash
curl -X DELETE http://localhost:8080/tasks/1
```

If something is wrong, the server sends back an error message in JSON, like `{"error":"Task not found"}`, with the right status code (400 for bad input, 404 if the task doesn't exist).

## Running the tests

I wrote some tests in `handler_test.go` using Go's `httptest` package, so they don't need the server to be running:

```bash
go test
```

## Project files

- `main.go` - starts the server on port 8080
- `server.go` - the `Server` struct that holds the tasks, plus the routes
- `handlers.go` - the code for each endpoint (create, list, get, delete)
- `task.go` - the `Task` struct
- `response.go` - a helper function for sending JSON responses
- `handler_test.go` - tests for the handlers

## What I learned

- How to use Go's `net/http` package to make a server without any frameworks
- How routing works and how to read values from the URL (like the task ID)
- How to read JSON from a request and send JSON back
- Why you need a mutex (`sync.Mutex`) when more than one request can change the same data at the same time
- Using HTTP status codes like 200, 201, 400, and 404 correctly
- Writing tests for HTTP handlers

## Things I might add later

- An endpoint to mark a task as completed (`PUT` or `PATCH`)
- Saving tasks to a file or database so they don't disappear
- Returning the task list sorted by ID

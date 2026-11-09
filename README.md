# HTTP Server (Task API)

This is a small HTTP server I built in Go to learn how web servers and REST APIs work. It lets you create, view, and delete tasks (like a really simple to-do list) by sending HTTP requests.

Tasks are saved in a PostgreSQL database, so they're still there after the server restarts.

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

You need:

- Go 1.22 or newer (the routes use the `"GET /tasks/{id}"` style patterns, which were added in 1.22)
- PostgreSQL ([download here](https://www.postgresql.org/download/))

1. Create a database for the app. In `psql` (or pgAdmin):

```sql
CREATE DATABASE tasks;
```

You don't need to create the table. The server creates the `tasks` table the first time it starts.

2. Clone the repo:

```bash
git clone https://github.com/AVargas1997/HTTP_Server.git
```

3. Go into the folder:

```bash
cd HTTP_Server
```

4. Tell the server how to connect to the database by setting `DATABASE_URL` (change `password` to your Postgres password).

PowerShell:

```powershell
$env:DATABASE_URL = "postgres://postgres:password@localhost:5432/tasks?sslmode=disable"
```

Mac/Linux:

```bash
export DATABASE_URL="postgres://postgres:password@localhost:5432/tasks?sslmode=disable"
```

5. Start the server:

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

The handler tests in `handler_test.go` use Go's `httptest` package and an in-memory store, so they don't need the server or a database:

```bash
go test
```

`postgres_store_test.go` tests the real database code. It only runs if `TEST_DATABASE_URL` is set, and it **empties the tasks table**, so use a separate database for it:

```sql
CREATE DATABASE tasks_test;
```

```powershell
$env:TEST_DATABASE_URL = "postgres://postgres:password@localhost:5432/tasks_test?sslmode=disable"
go test -v
```

## Project files

- `main.go` - connects to PostgreSQL and starts the server on port 8080
- `server.go` - the `Server` struct and the routes
- `store.go` - the `TaskStore` interface that both stores follow
- `postgres_store.go` - saves tasks in PostgreSQL using SQL queries
- `memory_store.go` - keeps tasks in a map (used by the tests)
- `handlers.go` - the code for each endpoint (create, list, get, delete)
- `task.go` - the `Task` struct
- `response.go` - a helper function for sending JSON responses
- `handler_test.go` - tests for the handlers
- `postgres_store_test.go` - tests for the PostgreSQL store

## What I learned

- How to use Go's `net/http` package to make a server without any frameworks
- How routing works and how to read values from the URL (like the task ID)
- How to read JSON from a request and send JSON back
- Why you need a mutex (`sync.Mutex`) when more than one request can change the same data at the same time
- Using HTTP status codes like 200, 201, 400, and 404 correctly
- Writing tests for HTTP handlers
- Connecting Go to PostgreSQL with `database/sql` and the `pgx` driver
- Using `$1` placeholders in SQL queries to prevent SQL injection
- Using an interface (`TaskStore`) so the handlers don't care where the data is saved

## Things I might add later

- An endpoint to mark a task as completed (`PUT` or `PATCH`)

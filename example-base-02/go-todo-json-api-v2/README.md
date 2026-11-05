# Go Todo JSON REST API – v2

Layered Go REST API using only the standard library and JSON file persistence.

## Requirements
- Go 1.22+

## Run
```powershell
go run ./cmd/server
```
API base URL: `http://localhost:8080`

Optional JSON path:
```powershell
$env:TODO_DATA_FILE = "data/todos.json"
go run ./cmd/server
```

## Structure
```text
go-todo-json-api-v2/
├── cmd/server/main.go
├── internal/
│   ├── handler/todo_handler.go
│   ├── model/todo.go
│   ├── repository/todo_repository.go
│   └── service/todo_service.go
├── data/todos.json
├── go.mod
└── README.md
```

## Endpoints

| Method | Endpoint | Description |
|---|---|---|
| GET | `/api/todos` | List todos |
| GET | `/api/todos/{id}` | Get one todo |
| POST | `/api/todos` | Create todo |
| PUT | `/api/todos/{id}` | Replace editable fields |
| PATCH | `/api/todos/{id}` | Update only supplied fields |
| DELETE | `/api/todos/{id}` | Delete todo |

### Create
```powershell
curl.exe -X POST http://localhost:8080/api/todos `
  -H "Content-Type: application/json" `
  -d '{\"title\":\"Learn Go\",\"description\":\"Build REST APIs\",\"completed\":false}'
```

### List all
```powershell
curl.exe http://localhost:8080/api/todos
```

### Filter by completion
```powershell
curl.exe "http://localhost:8080/api/todos?completed=false"
```

### Search title and description
```powershell
curl.exe "http://localhost:8080/api/todos?q=go"
```

### Get by ID
```powershell
curl.exe http://localhost:8080/api/todos/1
```

### Full update (PUT)
```powershell
curl.exe -X PUT http://localhost:8080/api/todos/1 `
  -H "Content-Type: application/json" `
  -d '{\"title\":\"Learn Go deeply\",\"description\":\"Practice concurrency\",\"completed\":true}'
```

### Partial update (PATCH)
```powershell
curl.exe -X PATCH http://localhost:8080/api/todos/1 `
  -H "Content-Type: application/json" `
  -d '{\"completed\":false}'
```

### Delete
```powershell
curl.exe -i -X DELETE http://localhost:8080/api/todos/1
```

## Tests and formatting
```powershell
gofmt -w .
go test ./...
go vet ./...
```

## v2 improvements
- Handler, service, repository and model layers.
- Repository interface for easier testing and future database replacement.
- PUT and PATCH support; PATCH can explicitly set booleans to false.
- Search (`q`) and completion (`completed`) filters.
- Input validation and consistent JSON errors.
- `sync.RWMutex` protects shared in-memory data.
- Graceful HTTP server shutdown.
- Unit tests for service behavior.

## Storage note
The JSON file is loaded at startup and rewritten for mutations. This is suitable for learning and a single application process, not concurrent multi-process production deployments.

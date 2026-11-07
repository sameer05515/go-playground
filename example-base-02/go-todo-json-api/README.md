# Go Todo JSON REST API

A REST API using Go's standard library. Data is persisted in `data/todos.json`.

## Requirements
Go 1.22 or later.

## Run
```powershell
go run .
```
Server: `http://localhost:8080`

## Endpoints

| Method | Endpoint | Description |
|---|---|---|
| GET | `/api/todos` | List todos |
| GET | `/api/todos/{id}` | Get a todo |
| POST | `/api/todos` | Create a todo |
| PUT | `/api/todos/{id}` | Update a todo |
| DELETE | `/api/todos/{id}` | Delete a todo |

## PowerShell examples

Create:
```powershell
curl.exe -X POST http://localhost:8080/api/todos `
  -H "Content-Type: application/json" `
  -d '{\"title\":\"Learn Go\",\"description\":\"Build a JSON REST API\",\"completed\":false}'
```

List:
```powershell
curl.exe http://localhost:8080/api/todos
```

Get:
```powershell
curl.exe http://localhost:8080/api/todos/1
```

Update (PUT replaces all editable fields):
```powershell
curl.exe -X PUT http://localhost:8080/api/todos/1 `
  -H "Content-Type: application/json" `
  -d '{\"title\":\"Learn Go REST APIs\",\"description\":\"Practice CRUD\",\"completed\":true}'
```

Delete:
```powershell
curl.exe -i -X DELETE http://localhost:8080/api/todos/1
```

## Design notes
- The JSON file is loaded at startup and rewritten after every successful mutation.
- `sync.RWMutex` protects in-memory data from concurrent HTTP requests.
- Writes go to a temporary file and are then renamed over the data file.
- This is a learning project. A JSON file is not suitable for multi-process or high-write production workloads.

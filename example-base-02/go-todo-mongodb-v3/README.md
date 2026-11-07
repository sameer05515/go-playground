# Go Todo MongoDB API – V3

Requires Go 1.23+, Docker Desktop, and the official MongoDB Go driver v2.

```powershell
docker compose up -d
go mod tidy
go run ./cmd/server
```

Defaults: `MONGO_URI=mongodb://localhost:27017`, `MONGO_DB=todo_db`, `PORT=:8080`.

| Method | Route | Action |
|---|---|---|
| GET | /api/todos | List, filter `?completed=true`, search `?q=go` |
| GET | /api/todos/{id} | Get by ObjectID |
| POST | /api/todos | Create |
| PUT | /api/todos/{id} | Replace editable fields |
| PATCH | /api/todos/{id} | Partial update |
| DELETE | /api/todos/{id} | Delete |

PowerShell examples:

```powershell
$body = @{title='Learn Go';description='MongoDB CRUD';completed=$false} | ConvertTo-Json
$todo = Invoke-RestMethod -Method Post -Uri http://localhost:8080/api/todos -ContentType application/json -Body $body
$todo
Invoke-RestMethod http://localhost:8080/api/todos
Invoke-RestMethod "http://localhost:8080/api/todos/$($todo.id)"
$patch = @{completed=$true} | ConvertTo-Json
Invoke-RestMethod -Method Patch -Uri "http://localhost:8080/api/todos/$($todo.id)" -ContentType application/json -Body $patch
Invoke-RestMethod -Method Delete -Uri "http://localhost:8080/api/todos/$($todo.id)"
```

MongoDB persists in Docker named volume `mongodb_data`. This demo runs MongoDB without authentication, bound to localhost only. For production enable authentication, TLS and appropriate network restrictions. MongoDB `_id` values are represented as hex strings in JSON.

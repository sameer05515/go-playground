package main

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Todo struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Completed   bool      `json:"completed"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type TodoInput struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Completed   bool   `json:"completed"`
}

type TodoStore struct {
	mu     sync.RWMutex
	todos  []Todo
	nextID int
	file   string
}

func NewTodoStore(file string) (*TodoStore, error) {
	s := &TodoStore{todos: []Todo{}, nextID: 1, file: file}
	if err := os.MkdirAll(filepath.Dir(file), 0755); err != nil {
		return nil, err
	}

	data, err := os.ReadFile(file)
	if errors.Is(err, os.ErrNotExist) {
		if err := s.persistLocked(); err != nil {
			return nil, err
		}
		return s, nil
	}
	if err != nil {
		return nil, err
	}
	if len(strings.TrimSpace(string(data))) > 0 {
		if err := json.Unmarshal(data, &s.todos); err != nil {
			return nil, err
		}
	}
	for _, todo := range s.todos {
		if todo.ID >= s.nextID {
			s.nextID = todo.ID + 1
		}
	}
	return s, nil
}

func (s *TodoStore) persistLocked() error {
	data, err := json.MarshalIndent(s.todos, "", "  ")
	if err != nil {
		return err
	}
	temp := s.file + ".tmp"
	if err := os.WriteFile(temp, data, 0644); err != nil {
		return err
	}
	return os.Rename(temp, s.file)
}

func (s *TodoStore) List() []Todo {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result := make([]Todo, len(s.todos))
	copy(result, s.todos)
	return result
}

func (s *TodoStore) Get(id int) (Todo, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, todo := range s.todos {
		if todo.ID == id {
			return todo, true
		}
	}
	return Todo{}, false
}

func (s *TodoStore) Create(input TodoInput) (Todo, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now().UTC()
	todo := Todo{
		ID: s.nextID, Title: strings.TrimSpace(input.Title),
		Description: strings.TrimSpace(input.Description),
		Completed: input.Completed, CreatedAt: now, UpdatedAt: now,
	}
	s.todos = append(s.todos, todo)
	s.nextID++
	if err := s.persistLocked(); err != nil {
		s.todos = s.todos[:len(s.todos)-1]
		s.nextID--
		return Todo{}, err
	}
	return todo, nil
}

func (s *TodoStore) Update(id int, input TodoInput) (Todo, bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i := range s.todos {
		if s.todos[i].ID == id {
			old := s.todos[i]
			s.todos[i].Title = strings.TrimSpace(input.Title)
			s.todos[i].Description = strings.TrimSpace(input.Description)
			s.todos[i].Completed = input.Completed
			s.todos[i].UpdatedAt = time.Now().UTC()
			if err := s.persistLocked(); err != nil {
				s.todos[i] = old
				return Todo{}, true, err
			}
			return s.todos[i], true, nil
		}
	}
	return Todo{}, false, nil
}

func (s *TodoStore) Delete(id int) (bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i := range s.todos {
		if s.todos[i].ID == id {
			removed := s.todos[i]
			s.todos = append(s.todos[:i], s.todos[i+1:]...)
			if err := s.persistLocked(); err != nil {
				s.todos = append(s.todos, Todo{})
				copy(s.todos[i+1:], s.todos[i:])
				s.todos[i] = removed
				return false, err
			}
			return true, nil
		}
	}
	return false, nil
}

type API struct{ store *TodoStore }

func main() {
	store, err := NewTodoStore("data/todos.json")
	if err != nil {
		log.Fatalf("initialize todo store: %v", err)
	}
	api := &API{store: store}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/todos", api.listTodos)
	mux.HandleFunc("POST /api/todos", api.createTodo)
	mux.HandleFunc("GET /api/todos/{id}", api.getTodo)
	mux.HandleFunc("PUT /api/todos/{id}", api.updateTodo)
	mux.HandleFunc("DELETE /api/todos/{id}", api.deleteTodo)

	server := &http.Server{
		Addr: ":8080", Handler: loggingMiddleware(mux),
		ReadHeaderTimeout: 5 * time.Second,
	}
	log.Println("Todo API running at http://localhost:8080")
	log.Fatal(server.ListenAndServe())
}

func (a *API) listTodos(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, a.store.List())
}

func (a *API) getTodo(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r)
	if !ok { return }
	todo, found := a.store.Get(id)
	if !found {
		writeError(w, http.StatusNotFound, "todo not found")
		return
	}
	writeJSON(w, http.StatusOK, todo)
}

func (a *API) createTodo(w http.ResponseWriter, r *http.Request) {
	var input TodoInput
	if !decodeJSON(w, r, &input) { return }
	if strings.TrimSpace(input.Title) == "" {
		writeError(w, http.StatusBadRequest, "title is required")
		return
	}
	todo, err := a.store.Create(input)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "could not save todo")
		return
	}
	writeJSON(w, http.StatusCreated, todo)
}

func (a *API) updateTodo(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r)
	if !ok { return }
	var input TodoInput
	if !decodeJSON(w, r, &input) { return }
	if strings.TrimSpace(input.Title) == "" {
		writeError(w, http.StatusBadRequest, "title is required")
		return
	}
	todo, found, err := a.store.Update(id, input)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "could not save todo")
		return
	}
	if !found {
		writeError(w, http.StatusNotFound, "todo not found")
		return
	}
	writeJSON(w, http.StatusOK, todo)
}

func (a *API) deleteTodo(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r)
	if !ok { return }
	deleted, err := a.store.Delete(id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "could not save todo")
		return
	}
	if !deleted {
		writeError(w, http.StatusNotFound, "todo not found")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func parseID(w http.ResponseWriter, r *http.Request) (int, bool) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id <= 0 {
		writeError(w, http.StatusBadRequest, "id must be a positive integer")
		return 0, false
	}
	return id, true
}

func decodeJSON(w http.ResponseWriter, r *http.Request, dst any) bool {
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(dst); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON: "+err.Error())
		return false
	}
	if err := decoder.Decode(&struct{}{}); err != io.EOF {
		writeError(w, http.StatusBadRequest, "request must contain exactly one JSON object")
		return false
	}
	return true
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}

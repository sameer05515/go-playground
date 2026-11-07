package repository

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"sync"

	"go-todo-json-api-v2/internal/model"
)

var ErrNotFound = errors.New("todo not found")

type TodoRepository interface {
	List() ([]model.Todo, error)
	GetByID(id int) (model.Todo, error)
	Create(todo model.Todo) (model.Todo, error)
	Update(todo model.Todo) (model.Todo, error)
	Delete(id int) error
}

type JSONTodoRepository struct {
	mu     sync.RWMutex
	todos  []model.Todo
	nextID int
	file   string
}

func NewJSONTodoRepository(file string) (*JSONTodoRepository, error) {
	repo := &JSONTodoRepository{
		todos:  make([]model.Todo, 0),
		nextID: 1,
		file:   file,
	}

	if err := os.MkdirAll(filepath.Dir(file), 0755); err != nil {
		return nil, err
	}

	data, err := os.ReadFile(file)
	if errors.Is(err, os.ErrNotExist) {
		if err := repo.persistLocked(); err != nil {
			return nil, err
		}
		return repo, nil
	}
	if err != nil {
		return nil, err
	}

	if len(data) > 0 {
		if err := json.Unmarshal(data, &repo.todos); err != nil {
			return nil, err
		}
	}
	if repo.todos == nil {
		repo.todos = make([]model.Todo, 0)
	}

	for _, todo := range repo.todos {
		if todo.ID >= repo.nextID {
			repo.nextID = todo.ID + 1
		}
	}
	return repo, nil
}

func (r *JSONTodoRepository) List() ([]model.Todo, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make([]model.Todo, len(r.todos))
	copy(result, r.todos)
	return result, nil
}

func (r *JSONTodoRepository) GetByID(id int) (model.Todo, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, todo := range r.todos {
		if todo.ID == id {
			return todo, nil
		}
	}
	return model.Todo{}, ErrNotFound
}

func (r *JSONTodoRepository) Create(todo model.Todo) (model.Todo, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	todo.ID = r.nextID
	r.todos = append(r.todos, todo)
	if err := r.persistLocked(); err != nil {
		r.todos = r.todos[:len(r.todos)-1]
		return model.Todo{}, err
	}
	r.nextID++
	return todo, nil
}

func (r *JSONTodoRepository) Update(todo model.Todo) (model.Todo, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	for i := range r.todos {
		if r.todos[i].ID == todo.ID {
			previous := r.todos[i]
			r.todos[i] = todo
			if err := r.persistLocked(); err != nil {
				r.todos[i] = previous
				return model.Todo{}, err
			}
			return todo, nil
		}
	}
	return model.Todo{}, ErrNotFound
}

func (r *JSONTodoRepository) Delete(id int) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	for i := range r.todos {
		if r.todos[i].ID == id {
			previous := r.todos[i]
			r.todos = append(r.todos[:i], r.todos[i+1:]...)
			if err := r.persistLocked(); err != nil {
				r.todos = append(r.todos, model.Todo{})
				copy(r.todos[i+1:], r.todos[i:])
				r.todos[i] = previous
				return err
			}
			return nil
		}
	}
	return ErrNotFound
}

func (r *JSONTodoRepository) persistLocked() error {
	data, err := json.MarshalIndent(r.todos, "", "  ")
	if err != nil {
		return err
	}

	tempFile := r.file + ".tmp"
	if err := os.WriteFile(tempFile, data, 0644); err != nil {
		return err
	}
	return os.Rename(tempFile, r.file)
}

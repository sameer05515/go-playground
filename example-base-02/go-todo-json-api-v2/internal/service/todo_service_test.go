package service

import (
	"errors"
	"testing"

	"go-todo-json-api-v2/internal/model"
	"go-todo-json-api-v2/internal/repository"
)

type memoryRepository struct {
	todos []model.Todo
	next  int
}

func (m *memoryRepository) List() ([]model.Todo, error) {
	result := append([]model.Todo(nil), m.todos...)
	return result, nil
}

func (m *memoryRepository) GetByID(id int) (model.Todo, error) {
	for _, todo := range m.todos {
		if todo.ID == id {
			return todo, nil
		}
	}
	return model.Todo{}, repository.ErrNotFound
}

func (m *memoryRepository) Create(todo model.Todo) (model.Todo, error) {
	todo.ID = m.next
	m.next++
	m.todos = append(m.todos, todo)
	return todo, nil
}

func (m *memoryRepository) Update(todo model.Todo) (model.Todo, error) {
	for i := range m.todos {
		if m.todos[i].ID == todo.ID {
			m.todos[i] = todo
			return todo, nil
		}
	}
	return model.Todo{}, repository.ErrNotFound
}

func (m *memoryRepository) Delete(id int) error {
	for i := range m.todos {
		if m.todos[i].ID == id {
			m.todos = append(m.todos[:i], m.todos[i+1:]...)
			return nil
		}
	}
	return repository.ErrNotFound
}

func TestCreateRejectsBlankTitle(t *testing.T) {
	svc := NewTodoService(&memoryRepository{next: 1})
	_, err := svc.Create(model.CreateTodoInput{Title: "   "})
	if !errors.Is(err, ErrInvalidTitle) {
		t.Fatalf("expected ErrInvalidTitle, got %v", err)
	}
}

func TestPatchCanSetCompletedFalse(t *testing.T) {
	repo := &memoryRepository{
		next: 2,
		todos: []model.Todo{{ID: 1, Title: "Test", Completed: true}},
	}
	svc := NewTodoService(repo)

	falseValue := false
	updated, err := svc.Patch(1, model.PatchTodoInput{Completed: &falseValue})
	if err != nil {
		t.Fatalf("Patch returned error: %v", err)
	}
	if updated.Completed {
		t.Fatal("expected Completed to be false after PATCH")
	}
}

func TestListFiltersByCompletedAndSearch(t *testing.T) {
	repo := &memoryRepository{
		next: 3,
		todos: []model.Todo{
			{ID: 1, Title: "Learn Go", Completed: true},
			{ID: 2, Title: "Learn Java", Completed: false},
		},
	}
	svc := NewTodoService(repo)

	completed := true
	result, err := svc.List(model.TodoFilter{Completed: &completed, Search: "go"})
	if err != nil {
		t.Fatalf("List returned error: %v", err)
	}
	if len(result) != 1 || result[0].ID != 1 {
		t.Fatalf("unexpected filtered result: %+v", result)
	}
}

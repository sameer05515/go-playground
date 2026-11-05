package service

import (
	"errors"
	"strings"
	"time"

	"go-todo-json-api-v2/internal/model"
	"go-todo-json-api-v2/internal/repository"
)

var ErrInvalidTitle = errors.New("title is required")

type TodoService struct {
	repo repository.TodoRepository
}

func NewTodoService(repo repository.TodoRepository) *TodoService {
	return &TodoService{repo: repo}
}

func (s *TodoService) List(filter model.TodoFilter) ([]model.Todo, error) {
	todos, err := s.repo.List()
	if err != nil {
		return nil, err
	}

	filtered := make([]model.Todo, 0, len(todos))
	search := strings.ToLower(strings.TrimSpace(filter.Search))

	for _, todo := range todos {
		if filter.Completed != nil && todo.Completed != *filter.Completed {
			continue
		}
		if search != "" &&
			!strings.Contains(strings.ToLower(todo.Title), search) &&
			!strings.Contains(strings.ToLower(todo.Description), search) {
			continue
		}
		filtered = append(filtered, todo)
	}
	return filtered, nil
}

func (s *TodoService) GetByID(id int) (model.Todo, error) {
	return s.repo.GetByID(id)
}

func (s *TodoService) Create(input model.CreateTodoInput) (model.Todo, error) {
	title := strings.TrimSpace(input.Title)
	if title == "" {
		return model.Todo{}, ErrInvalidTitle
	}

	now := time.Now().UTC()
	return s.repo.Create(model.Todo{
		Title:       title,
		Description: strings.TrimSpace(input.Description),
		Completed:   input.Completed,
		CreatedAt:   now,
		UpdatedAt:   now,
	})
}

func (s *TodoService) Update(id int, input model.UpdateTodoInput) (model.Todo, error) {
	title := strings.TrimSpace(input.Title)
	if title == "" {
		return model.Todo{}, ErrInvalidTitle
	}

	current, err := s.repo.GetByID(id)
	if err != nil {
		return model.Todo{}, err
	}

	current.Title = title
	current.Description = strings.TrimSpace(input.Description)
	current.Completed = input.Completed
	current.UpdatedAt = time.Now().UTC()
	return s.repo.Update(current)
}

func (s *TodoService) Patch(id int, input model.PatchTodoInput) (model.Todo, error) {
	current, err := s.repo.GetByID(id)
	if err != nil {
		return model.Todo{}, err
	}

	if input.Title != nil {
		title := strings.TrimSpace(*input.Title)
		if title == "" {
			return model.Todo{}, ErrInvalidTitle
		}
		current.Title = title
	}
	if input.Description != nil {
		current.Description = strings.TrimSpace(*input.Description)
	}
	if input.Completed != nil {
		current.Completed = *input.Completed
	}

	current.UpdatedAt = time.Now().UTC()
	return s.repo.Update(current)
}

func (s *TodoService) Delete(id int) error {
	return s.repo.Delete(id)
}

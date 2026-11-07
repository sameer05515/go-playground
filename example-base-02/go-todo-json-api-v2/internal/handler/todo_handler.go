package handler

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"
	"strings"

	"go-todo-json-api-v2/internal/model"
	"go-todo-json-api-v2/internal/repository"
	"go-todo-json-api-v2/internal/service"
)

type TodoHandler struct {
	service *service.TodoService
}

func NewTodoHandler(todoService *service.TodoService) *TodoHandler {
	return &TodoHandler{service: todoService}
}

func (h *TodoHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/todos", h.list)
	mux.HandleFunc("POST /api/todos", h.create)
	mux.HandleFunc("GET /api/todos/{id}", h.get)
	mux.HandleFunc("PUT /api/todos/{id}", h.update)
	mux.HandleFunc("PATCH /api/todos/{id}", h.patch)
	mux.HandleFunc("DELETE /api/todos/{id}", h.delete)
}

func (h *TodoHandler) list(w http.ResponseWriter, r *http.Request) {
	filter := model.TodoFilter{Search: r.URL.Query().Get("q")}
	if raw := r.URL.Query().Get("completed"); raw != "" {
		value, err := strconv.ParseBool(raw)
		if err != nil {
			writeError(w, http.StatusBadRequest, "completed must be true or false")
			return
		}
		filter.Completed = &value
	}

	todos, err := h.service.List(filter)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "could not list todos")
		return
	}
	writeJSON(w, http.StatusOK, todos)
}

func (h *TodoHandler) get(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(w, r)
	if !ok {
		return
	}
	todo, err := h.service.GetByID(id)
	if err != nil {
		handleServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, todo)
}

func (h *TodoHandler) create(w http.ResponseWriter, r *http.Request) {
	var input model.CreateTodoInput
	if !decodeJSON(w, r, &input) {
		return
	}
	todo, err := h.service.Create(input)
	if err != nil {
		handleServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, todo)
}

func (h *TodoHandler) update(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(w, r)
	if !ok {
		return
	}
	var input model.UpdateTodoInput
	if !decodeJSON(w, r, &input) {
		return
	}
	todo, err := h.service.Update(id, input)
	if err != nil {
		handleServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, todo)
}

func (h *TodoHandler) patch(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(w, r)
	if !ok {
		return
	}
	var input model.PatchTodoInput
	if !decodeJSON(w, r, &input) {
		return
	}
	todo, err := h.service.Patch(id, input)
	if err != nil {
		handleServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, todo)
}

func (h *TodoHandler) delete(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(w, r)
	if !ok {
		return
	}
	if err := h.service.Delete(id); err != nil {
		handleServiceError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func pathID(w http.ResponseWriter, r *http.Request) (int, bool) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id < 1 {
		writeError(w, http.StatusBadRequest, "id must be a positive integer")
		return 0, false
	}
	return id, true
}

func decodeJSON(w http.ResponseWriter, r *http.Request, destination any) bool {
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(destination); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON: "+err.Error())
		return false
	}
	if err := decoder.Decode(&struct{}{}); err != io.EOF {
		writeError(w, http.StatusBadRequest, "request must contain exactly one JSON object")
		return false
	}
	return true
}

func handleServiceError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, service.ErrInvalidTitle):
		writeError(w, http.StatusBadRequest, err.Error())
	case errors.Is(err, repository.ErrNotFound):
		writeError(w, http.StatusNotFound, err.Error())
	default:
		writeError(w, http.StatusInternalServerError, "internal server error")
	}
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": strings.TrimSpace(message)})
}

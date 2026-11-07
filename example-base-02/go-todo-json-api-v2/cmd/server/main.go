package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go-todo-json-api-v2/internal/handler"
	"go-todo-json-api-v2/internal/repository"
	"go-todo-json-api-v2/internal/service"
)

func main() {
	dataFile := os.Getenv("TODO_DATA_FILE")
	if dataFile == "" {
		dataFile = "data/todos.json"
	}

	repo, err := repository.NewJSONTodoRepository(dataFile)
	if err != nil {
		log.Fatalf("initialize JSON repository: %v", err)
	}

	todoService := service.NewTodoService(repo)
	todoHandler := handler.NewTodoHandler(todoService)

	mux := http.NewServeMux()
	todoHandler.RegisterRoutes(mux)

	server := &http.Server{
		Addr:              ":8080",
		Handler:           requestLogger(mux),
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		log.Println("Todo API v2 listening on http://localhost:8080")
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("HTTP server: %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	log.Println("Shutting down server...")
	if err := server.Shutdown(ctx); err != nil {
		log.Printf("graceful shutdown failed: %v", err)
	}
}

func requestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("%s %s %s", r.Method, r.URL.Path, time.Since(start))
	})
}

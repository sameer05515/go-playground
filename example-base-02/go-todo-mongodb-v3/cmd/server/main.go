package main

import (
	"context"
	"errors"
	"go-todo-mongodb-v3/internal/handler"
	"go-todo-mongodb-v3/internal/repository"
	"go-todo-mongodb-v3/internal/service"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func env(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
func main() {
	uri := env("MONGO_URI", "mongodb://localhost:27017")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := client.Disconnect(context.Background()); err != nil {
			log.Print(err)
		}
	}()
	if err := client.Ping(ctx, nil); err != nil {
		log.Fatalf("MongoDB unavailable: %v", err)
	}
	repo := repository.New(client.Database(env("MONGO_DB", "todo_db")))
	h := handler.New(service.New(repo))
	mux := http.NewServeMux()
	h.Register(mux)
	server := &http.Server{Addr: env("PORT", ":8080"), Handler: mux, ReadHeaderTimeout: 5 * time.Second}
	go func() {
		log.Printf("API listening at %s", server.Addr)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal(err)
		}
	}()
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop
	shutdownCtx, done := context.WithTimeout(context.Background(), 10*time.Second)
	defer done()
	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Print(err)
	}
}

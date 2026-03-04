package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"issue-tracker/backend/internal/api"
	"issue-tracker/backend/internal/api/handler"
	"issue-tracker/backend/internal/store/mongodb"
)

func main() {
	ctx := context.Background()
	client, err := mongodb.NewClient(ctx)
	if err != nil {
		log.Fatalf("MongoDB connect: %v", err)
	}
	defer func() {
		if err := client.Disconnect(ctx); err != nil {
			log.Printf("MongoDB disconnect: %v", err)
		}
	}()

	db := mongodb.Database(client)
	router := api.NewRouter(db, log.Default())
	router.Handle("/issues", handler.HandlePostIssues(router))
	router.Handle("/users", handler.HandleGetUsers(router))

	server := &http.Server{Addr: ":8080", Handler: router}
	go func() {
		log.Println("Server listening on :8080")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	_ = server.Shutdown(context.Background())
}

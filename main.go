package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

const storageFile = "url_mappings.json"

func main() {
	loadMappings()

	http.HandleFunc("/shorten", shortenURLHandler)
	http.HandleFunc("/", redirectHandler)
	srv := &http.Server{Addr: ":8080"}

	go func() {
		log.Println("Starting server on :8080")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	saveMappings()

	log.Println("Server exiting")
}

package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"urlshortener/urlshortener"
)

func main() {
	urlshortener.LoadMappings(urlshortener.StorageFile)

	http.HandleFunc("/shorten", urlshortener.ShortenURLHandler)
	http.HandleFunc("/", urlshortener.RedirectHandler)
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

	urlshortener.SaveMappings(urlshortener.StorageFile)

	log.Println("Server exiting")
}

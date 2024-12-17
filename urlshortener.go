package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"os"
	"sync"
	"time"
)

var (
	shortenedURLs = make(map[string]string)
	urlMapMutex   sync.RWMutex
)

var src = rand.NewSource(time.Now().UnixNano())
var r = rand.New(src)

func generateShortURL() string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, 6)
	for i := range b {
		b[i] = letters[r.Intn(len(letters))]
	}
	return string(b)
}

func shortenURLHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		URL string `json:"url"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload\n"+err.Error(), http.StatusBadRequest)
		return
	}

	shortURL := generateShortURL()
	urlMapMutex.Lock()
	shortenedURLs[shortURL] = req.URL
	urlMapMutex.Unlock()

	saveMappings()

	resp := map[string]string{"short_url": "http://localhost:8080/" + shortURL}
	json.NewEncoder(w).Encode(resp)
}

func redirectHandler(w http.ResponseWriter, r *http.Request) {
	shortURL := r.URL.Path[1:]

	urlMapMutex.RLock()
	originalURL, ok := shortenedURLs[shortURL]
	urlMapMutex.RUnlock()

	if !ok {
		http.NotFound(w, r)
		return
	}

	http.Redirect(w, r, originalURL, http.StatusFound)
}

func loadMappings() {
	file, err := os.Open(storageFile)
	if err != nil {
		if os.IsNotExist(err) {
			return
		}
		log.Fatalf("Failed to open storage file: %v", err)
	}
	defer file.Close()

	urlMapMutex.Lock()
	defer urlMapMutex.Unlock()
	if err := json.NewDecoder(file).Decode(&shortenedURLs); err != nil {
		log.Fatalf("Failed to decode storage file: %v", err)
	}
}

func saveMappings() {
	file, err := os.Create(storageFile)
	if err != nil {
		log.Fatalf("Failed to create storage file: %v", err)
	}
	defer file.Close()

	urlMapMutex.RLock()
	defer urlMapMutex.RUnlock()
	if err := json.NewEncoder(file).Encode(shortenedURLs); err != nil {
		log.Fatalf("Failed to encode storage file: %v", err)
	}
}

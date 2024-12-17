package urlshortener

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"sync"
	"time"
)

const StorageFile = "url_mappings.json"

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

func isValidURL(u string) bool {
	_, err := url.ParseRequestURI(u)
	return err == nil
}

func ShortenURLHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		URL string `json:"url"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload\n"+err.Error(), http.StatusBadRequest)
		return
	}

	if !isValidURL(req.URL) {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}

	shortURL := generateShortURL()
	urlMapMutex.Lock()
	shortenedURLs[shortURL] = req.URL
	urlMapMutex.Unlock()

	SaveMappings(StorageFile)

	resp := map[string]string{"short_url": "http://localhost:8080/" + shortURL}
	json.NewEncoder(w).Encode(resp)
}

func RedirectHandler(w http.ResponseWriter, r *http.Request) {
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

func LoadMappings(filePath string) {
	file, err := os.Open(filePath)
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

func SaveMappings(filePath string) {
	file, err := os.Create(filePath)
	if err != nil {
		log.Fatalf("Failed to create storage file: %v", err)
	}
	defer file.Close()

	urlMapMutex.RLock()
	defer urlMapMutex.RUnlock()
	if err := json.NewEncoder(file).Encode(&shortenedURLs); err != nil {
		log.Fatalf("Failed to encode storage file: %v", err)
	}
}

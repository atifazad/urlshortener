package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestShortenURLHandler(t *testing.T) {
	// Create a new HTTP request with a JSON payload
	payload := map[string]string{"url": "https://www.example.com"}
	payloadBytes, _ := json.Marshal(payload)
	req, err := http.NewRequest("POST", "/shorten", bytes.NewBuffer(payloadBytes))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Create a ResponseRecorder to capture the response
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(shortenURLHandler)

	// Call the handler
	handler.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Check the response body
	var resp map[string]string
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}
	if _, ok := resp["short_url"]; !ok {
		t.Errorf("Response does not contain short_url")
	}
}

func TestRedirectHandler(t *testing.T) {
	// Add a test URL mapping
	shortenedURLs["test123"] = "https://www.example.com"

	// Create a new HTTP request
	req, err := http.NewRequest("GET", "/test123", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	// Create a ResponseRecorder to capture the response
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(redirectHandler)

	// Call the handler
	handler.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusFound {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusFound)
	}

	// Check the Location header
	expectedLocation := "https://www.example.com"
	if location := rr.Header().Get("Location"); location != expectedLocation {
		t.Errorf("handler returned wrong location header: got %v want %v", location, expectedLocation)
	}
}

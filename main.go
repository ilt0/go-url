package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strings"
	
)

type ShortenedURL struct {
	OriginalURL string `json:"original_url"`
	ShortURL    string `json:"short_url"`
}

var shortenedURLs map[string]string

func generateRandomKey(length int) string {
	const chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, length)
	for i := range result {
		result[i] = chars[rand.Intn(len(chars))]
	}
	return string(result)
}

func shortenURLHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	decoder := json.NewDecoder(r.Body)
	var requestBody map[string]string
	err := decoder.Decode(&requestBody)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	originalURL, ok := requestBody["url"]
	if !ok || originalURL == "" {
		http.Error(w, "Missing or empty 'url' field in the request body", http.StatusBadRequest)
		return
	}

	shortURLKey := generateRandomKey(6)
	shortenedURLs[shortURLKey] = originalURL

	shortenedURL := ShortenedURL{
		OriginalURL: originalURL,
		ShortURL:    fmt.Sprintf("/%s", shortURLKey),
	}

	responseJSON, err := json.Marshal(shortenedURL)
	if err != nil {
		http.Error(w, "Failed to marshal JSON response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(responseJSON)
}

func redirectHandler(w http.ResponseWriter, r *http.Request) {
	shortURLKey := strings.TrimPrefix(r.URL.Path, "/")
	originalURL, found := shortenedURLs[shortURLKey]
	if !found {
		http.NotFound(w, r)
		return
	}

	http.Redirect(w, r, originalURL, http.StatusSeeOther)
}

func serveIndexPage(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "index.html")
}

func main() {
	shortenedURLs = make(map[string]string)
	rand.Seed(42)

	http.HandleFunc("/shorten", shortenURLHandler)
	http.HandleFunc("/", redirectHandler)
	http.HandleFunc("/home", serveIndexPage)

	fmt.Println("Server started at http://localhost:3000")
	log.Fatal(http.ListenAndServe(":3000", nil))
}

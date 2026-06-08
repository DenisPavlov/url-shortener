package main

import (
	"crypto/sha256"
	"fmt"
	"io"
	"net/http"
)

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

// in memory storage
var urls = make(map[string]string)

func run() error {
	mux := http.NewServeMux()
	mux.HandleFunc("/", shortHandler)
	mux.HandleFunc("/{hash}", getUrlHandler)

	return http.ListenAndServe(`:8080`, mux)
}

func getUrlHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	hash := r.PathValue("hash")

	url, ok := urls[hash]
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.Header().Set("Location", url)
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func shortHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if r.Header.Get("Content-Type") != "text/plain" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	url := string(body)

	hash, err := hasher(url)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	urls[hash] = url

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)

	host := r.Host
	_, err = w.Write([]byte("http://" + host + "/" + hash))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

// todo -переделать функцию кеширования
func hasher(url string) (string, error) {
	h := sha256.New()
	if _, err := io.WriteString(h, url); err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", h.Sum(nil)), nil
}

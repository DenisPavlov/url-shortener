package main

import (
	"crypto/sha256"
	"fmt"
	"io"
	"net/http"

	"github.com/DenisPavlov/url-shortener/internal/handler"
	"github.com/go-chi/chi/v5"
)

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

// in memory storage
var urls = make(map[string]string)

func run() error {
	r := NewRouter()
	return http.ListenAndServe(`:8080`, r)
}

func NewRouter() chi.Router {
	r := chi.NewRouter()
	handler.Add(r, urls, hasher)
	return r
}

// todo -переделать функцию кеширования
func hasher(url string) (string, error) {
	h := sha256.New()
	if _, err := io.WriteString(h, url); err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", h.Sum(nil)), nil
}

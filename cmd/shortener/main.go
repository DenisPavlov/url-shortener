package main

import (
	"crypto/sha256"
	"fmt"
	"io"
	"net/http"

	"github.com/DenisPavlov/url-shortener/internal/handler"
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
	mux.HandleFunc("/", handler.Short(urls, hasher))
	mux.HandleFunc("/{hash}", handler.GetUrl(urls))

	return http.ListenAndServe(`:8080`, mux)
}

// todo -переделать функцию кеширования
func hasher(url string) (string, error) {
	h := sha256.New()
	if _, err := io.WriteString(h, url); err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", h.Sum(nil)), nil
}

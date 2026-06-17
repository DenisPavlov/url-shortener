package main

import (
	"crypto/sha256"
	"fmt"
	"io"
	"net/http"

	"github.com/DenisPavlov/url-shortener/internal/config"
	"github.com/DenisPavlov/url-shortener/internal/handler"
	"github.com/go-chi/chi/v5"
)

func main() {
	fmt.Println("Starting server...")

	cfg := config.MustLoad()
	fmt.Println("Config:", cfg)

	if err := run(cfg); err != nil {
		panic(err)
	}
}

// in memory storage
var urls = make(map[string]string)

func run(cfg config.Cfg) error {
	r := NewRouter(cfg)
	return http.ListenAndServe(cfg.ServerAddress, r)
}

func NewRouter(cfg config.Cfg) chi.Router {
	r := chi.NewRouter()
	handler.Add(r, urls, cfg.BaseURL, hasher)
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

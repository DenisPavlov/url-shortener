package main

import (
	"crypto/sha256"
	"fmt"
	"io"
	"net/http"

	"github.com/DenisPavlov/url-shortener/internal/config"
	"github.com/DenisPavlov/url-shortener/internal/handler"
	"github.com/DenisPavlov/url-shortener/internal/logger"
	"github.com/DenisPavlov/url-shortener/internal/middleware"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

func main() {
	fmt.Println("Starting server...")

	cfg := config.MustLoad()
	logger.MustInitialize(cfg.LogLevel)
	logger.Log.Info("Config was loaded successfully", zap.Any("config", cfg))

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
	r.Use(logger.RequestLogger)
	r.Use(middleware.Gzip)
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

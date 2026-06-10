package handler

import (
	"io"
	"net/http"
)

func Short(urls map[string]string, hasher func(string) (string, error)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "incorrect method", http.StatusBadRequest)
			return
		}
		if r.Header.Get("Content-Type") != "text/plain" {
			http.Error(w, "incorrect content type", http.StatusBadRequest)
			return
		}

		body, err := io.ReadAll(r.Body)
		if err != nil || len(body) == 0 {
			http.Error(w, "body must be not empty", http.StatusBadRequest)
			return
		}

		url := string(body)

		hash, err := hasher(url)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		urls[hash] = url

		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusCreated)

		host := r.Host
		_, err = io.WriteString(w, "http://"+host+"/"+hash)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
	}
}

func GetUrl(urls map[string]string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "incorrect method", http.StatusBadRequest)
			return
		}

		hash := r.PathValue("hash")

		url, ok := urls[hash]
		if !ok {
			http.Error(w, "url not found", http.StatusNotFound)
			return
		}

		w.Header().Set("Location", url)
		w.WriteHeader(http.StatusTemporaryRedirect)
	}
}

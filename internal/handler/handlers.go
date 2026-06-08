package handler

import (
	"io"
	"net/http"
)

func Short(urls map[string]string, hasher func(string) (string, error)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if r.Header.Get("Content-Type") != "text/plain" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		body, err := io.ReadAll(r.Body)
		if err != nil || len(body) == 0 {
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
}

func GetUrl(urls map[string]string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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
}

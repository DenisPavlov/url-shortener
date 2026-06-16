package handler

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"resty.dev/v3"
)

func TestShort(t *testing.T) {
	tests := []struct {
		name        string
		method      string
		contentType string
		url         string
		code        int
		isError     bool
	}{
		{
			name:        "happy path",
			method:      http.MethodPost,
			contentType: "text/plain",
			url:         "https://test.url",
			code:        http.StatusCreated,
			isError:     false,
		},
		{
			name:        "incorrect method",
			method:      http.MethodGet,
			contentType: "text/plain",
			code:        http.StatusMethodNotAllowed,
			isError:     true,
		},
		{
			name:    "unset Content-Type",
			method:  http.MethodPost,
			code:    http.StatusBadRequest,
			isError: true,
		},
		{
			name:        "empty body",
			method:      http.MethodPost,
			contentType: "text/plain",
			code:        http.StatusBadRequest,
			isError:     true,
		},
	}

	client := resty.New()
	client.SetRedirectPolicy(resty.RedirectNoPolicy())
	defer func(client *resty.Client) {
		_ = client.Close()
	}(client)

	urls := make(map[string]string)
	hasher := func(url string) (string, error) {
		return "aaa", nil
	}

	r := chi.NewRouter()
	Add(r, urls, hasher)
	srv := httptest.NewServer(r)
	defer srv.Close()

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {

			req := client.R()
			req.Method = tt.method
			req.SetContentType(tt.contentType)
			req.URL = fmt.Sprintf("%s/", srv.URL)
			req.Body = io.NopCloser(strings.NewReader(tt.url))

			resp, err := req.Send()
			require.NoError(t, err)

			assert.Equal(t, tt.code, resp.StatusCode())

			if tt.isError {
				return
			}

			assert.Equal(t, "text/plain", resp.Header().Get("Content-Type"))

			defer func(Body io.ReadCloser) {
				_ = Body.Close()
			}(resp.Body)
			resBody, err := io.ReadAll(resp.Body)
			require.NoError(t, err)
			assert.Contains(t, string(resBody), "aaa")
		})
	}
}

func TestGetUrls(t *testing.T) {
	tests := []struct {
		name    string
		method  string
		hash    string
		code    int
		isError bool
		hashes  map[string]string
	}{
		{
			name:    "happy path",
			method:  http.MethodGet,
			hash:    "aaa",
			code:    http.StatusTemporaryRedirect,
			isError: false,
			hashes: map[string]string{
				"aaa": "http://test.url",
			},
		},
		{
			name:    "incorrect method",
			method:  http.MethodPost,
			code:    http.StatusBadRequest,
			isError: true,
		},
		{
			name:    "hash not found",
			method:  http.MethodGet,
			hash:    "bbb",
			code:    http.StatusNotFound,
			isError: true,
			hashes:  make(map[string]string),
		},
	}

	client := resty.New()
	client.SetRedirectPolicy(resty.RedirectNoPolicy())
	defer func(client *resty.Client) {
		_ = client.Close()
	}(client)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := chi.NewRouter()
			Add(r, tt.hashes, nil)

			srv := httptest.NewServer(r)
			defer srv.Close()

			req := client.R()
			req.Method = tt.method
			req.URL = fmt.Sprintf("%s/%s", srv.URL, tt.hash)

			resp, err := req.Send()
			require.NoError(t, err)

			assert.Equal(t, tt.code, resp.StatusCode())
			if tt.isError {
				return
			}

			assert.Equal(t, "http://test.url", resp.Header().Get("Location"))

		})
	}
}

package handler

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
			code:        http.StatusBadRequest,
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
	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			urls := make(map[string]string)
			hasher := func(url string) (string, error) {
				return "aaa", nil
			}

			recorder := httptest.NewRecorder()
			req := httptest.NewRequest(tt.method, "/", strings.NewReader(tt.url))
			req.Header.Set("Content-Type", tt.contentType)

			h := Short(urls, hasher)
			h(recorder, req)

			resp := recorder.Result()

			assert.Equal(t, tt.code, resp.StatusCode)

			if tt.isError {
				return
			}

			assert.Equal(t, "text/plain", resp.Header.Get("Content-Type"))

			defer resp.Body.Close()
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
			code:    http.StatusNotFound,
			isError: true,
			hashes:  make(map[string]string),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			recorder := httptest.NewRecorder()
			req := httptest.NewRequest(tt.method, "/", strings.NewReader(tt.hash))
			req.SetPathValue("hash", tt.hash)

			h := GetUrl(tt.hashes)
			h(recorder, req)

			resp := recorder.Result()
			assert.Equal(t, tt.code, resp.StatusCode)
			if tt.isError {
				return
			}

			assert.Equal(t, "http://test.url", resp.Header.Get("Location"))

		})
	}
}

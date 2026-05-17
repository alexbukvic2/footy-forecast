package handler_test

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

// silentLogger returns a logger that discards all output. Use in tests where
// log lines would be noise.
func silentLogger() *slog.Logger {
	return slog.New(
		slog.NewTextHandler(
			&bytes.Buffer{}, &slog.HandlerOptions{
				Level: slog.LevelError + 1, // higher than Error → nothing logs
			},
		),
	)
}

// postJSON marshals body as JSON, posts it to the handler, returns the recorder.
func postJSON(
	t *testing.T,
	h http.HandlerFunc,
	path string,
	body any,
) *httptest.ResponseRecorder {
	t.Helper()
	raw, err := json.Marshal(body)
	require.NoError(t, err)
	req := httptest.NewRequest(http.MethodPost, path, bytes.NewReader(raw))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	h(rec, req)
	return rec
}

// postRaw posts a raw body string (used for testing malformed JSON).
func postRaw(
	t *testing.T,
	h http.HandlerFunc,
	path, body string,
) *httptest.ResponseRecorder {
	t.Helper()
	req := httptest.NewRequest(http.MethodPost, path, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	h(rec, req)
	return rec
}

// getWithPathValue calls a handler with a path-value parameter set
// (httptest.NewRequest does not route, so PathValue must be set manually).
func getWithPathValue(
	t *testing.T,
	h http.HandlerFunc,
	path, key, value string,
) *httptest.ResponseRecorder {
	t.Helper()
	req := httptest.NewRequest(http.MethodGet, path, nil)
	req.SetPathValue(key, value)
	rec := httptest.NewRecorder()
	h(rec, req)
	return rec
}

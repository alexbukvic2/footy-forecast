package handler_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/alexbukvic2/footy-forecast/internal/server/handler"
)

func TestSpec_ServeJSON(t *testing.T) {
	h := handler.NewSpec()
	rec := httptest.NewRecorder()
	h.ServeJSON(rec, httptest.NewRequest(http.MethodGet, "/openapi.json", nil))

	if rec.Code != http.StatusOK {
		t.Fatalf("want 200, got %d", rec.Code)
	}
	if ct := rec.Header().Get("Content-Type"); ct != "application/json" {
		t.Fatalf("want Content-Type application/json, got %q", ct)
	}
	if rec.Body.Len() == 0 {
		t.Fatal("expected non-empty body")
	}
}

func TestSpec_ServeUI(t *testing.T) {
	h := handler.NewSpec()
	rec := httptest.NewRecorder()
	h.ServeUI(rec, httptest.NewRequest(http.MethodGet, "/docs", nil))

	if rec.Code != http.StatusOK {
		t.Fatalf("want 200, got %d", rec.Code)
	}
	if ct := rec.Header().Get("Content-Type"); ct != "text/html; charset=utf-8" {
		t.Fatalf("want Content-Type text/html; charset=utf-8, got %q", ct)
	}
	if rec.Body.Len() == 0 {
		t.Fatal("expected non-empty body")
	}
}

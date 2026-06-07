package expo

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestClient_Send_Empty(t *testing.T) {
	t.Parallel()
	c := NewClient()
	receipts, err := c.Send(context.Background(), nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(receipts) != 0 {
		t.Fatalf("expected empty receipts, got %d", len(receipts))
	}
}

func TestClient_Send_HappyPath(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var msgs []Message
		if err := json.NewDecoder(r.Body).Decode(&msgs); err != nil {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}
		receipts := make([]Receipt, len(msgs))
		for i := range receipts {
			receipts[i] = Receipt{Status: "ok"}
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(sendResponse{Data: receipts})
	}))
	defer server.Close()

	c := NewClientWithEndpoint(server.URL, server.Client())

	msgs := []Message{
		{To: "ExponentPushToken[abc]", Title: "Hello", Body: "World"},
		{To: "ExponentPushToken[def]", Title: "Foo", Body: "Bar"},
	}

	receipts, err := c.Send(context.Background(), msgs)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(receipts) != 2 {
		t.Fatalf("expected 2 receipts, got %d", len(receipts))
	}
	for i, r := range receipts {
		if r.Status != "ok" {
			t.Errorf("receipt[%d].Status = %q, want %q", i, r.Status, "ok")
		}
	}
}

func TestClient_Send_Batching(t *testing.T) {
	t.Parallel()

	callCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callCount++
		var msgs []Message
		_ = json.NewDecoder(r.Body).Decode(&msgs)
		receipts := make([]Receipt, len(msgs))
		for i := range receipts {
			receipts[i] = Receipt{Status: "ok"}
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(sendResponse{Data: receipts})
	}))
	defer server.Close()

	c := NewClientWithEndpoint(server.URL, server.Client())

	// Send 150 messages — should create 2 batches.
	msgs := make([]Message, 150)
	for i := range msgs {
		msgs[i] = Message{To: "ExponentPushToken[x]", Title: "T", Body: "B"}
	}

	receipts, err := c.Send(context.Background(), msgs)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(receipts) != 150 {
		t.Fatalf("expected 150 receipts, got %d", len(receipts))
	}
	if callCount != 2 {
		t.Errorf("expected 2 HTTP calls for batching, got %d", callCount)
	}
}

func TestClient_Send_DeviceNotRegistered(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		receipts := []Receipt{{Status: "error", Details: &ReceiptDetails{Error: ErrDeviceNotRegistered}}}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(sendResponse{Data: receipts})
	}))
	defer server.Close()

	c := NewClientWithEndpoint(server.URL, server.Client())

	msgs := []Message{{To: "ExponentPushToken[stale]", Title: "T", Body: "B"}}
	receipts, err := c.Send(context.Background(), msgs)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(receipts) != 1 {
		t.Fatalf("expected 1 receipt, got %d", len(receipts))
	}
	if receipts[0].Status != "error" {
		t.Errorf("expected error status, got %q", receipts[0].Status)
	}
	if receipts[0].Details == nil || receipts[0].Details.Error != ErrDeviceNotRegistered {
		t.Errorf("expected DeviceNotRegistered in details, got %+v", receipts[0].Details)
	}
}

func TestClient_Send_TransientRetry(t *testing.T) {
	t.Parallel()

	callCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callCount++
		if callCount == 1 {
			// First call fails.
			http.Error(w, "service unavailable", http.StatusServiceUnavailable)
			return
		}
		// Second call succeeds.
		var msgs []Message
		_ = json.NewDecoder(r.Body).Decode(&msgs)
		receipts := []Receipt{{Status: "ok"}}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(sendResponse{Data: receipts})
	}))
	defer server.Close()

	c := NewClientWithEndpoint(server.URL, server.Client())

	msgs := []Message{{To: "ExponentPushToken[x]", Title: "T", Body: "B"}}
	receipts, err := c.Send(context.Background(), msgs)
	if err != nil {
		t.Fatalf("unexpected error after retry: %v", err)
	}
	if len(receipts) != 1 || receipts[0].Status != "ok" {
		t.Errorf("expected ok receipt after retry, got %+v", receipts)
	}
	if callCount != 2 {
		t.Errorf("expected 2 HTTP calls (1 fail + 1 retry), got %d", callCount)
	}
}

func TestClient_Send_TransientRetryExhausted(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		http.Error(w, "service unavailable", http.StatusServiceUnavailable)
	}))
	defer server.Close()

	c := NewClientWithEndpoint(server.URL, server.Client())

	msgs := []Message{{To: "ExponentPushToken[x]", Title: "T", Body: "B"}}
	_, err := c.Send(context.Background(), msgs)
	if err == nil {
		t.Fatal("expected error after exhausted retries, got nil")
	}
}

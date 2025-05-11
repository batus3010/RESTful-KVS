package kvs

import (
	"errors"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestInMemoryKVSIntegration(t *testing.T) {
	store := NewInMemoryKVS()

	// 1) Getting a missing key should error
	if val, err := store.Get("missing"); err == nil {
		t.Fatalf("expected error for missing key, got value %q and nil error", val)
	}

	// 2) Put a key, then Get should return it
	if err := store.Put("foo", "bar"); err != nil {
		t.Fatalf("unexpected error on Put: %v", err)
	}
	got, err := store.Get("foo")
	if err != nil {
		t.Fatalf("unexpected error on Get after Put: %v", err)
	}
	if got != "bar" {
		t.Fatalf("Get returned %q, want %q", got, "bar")
	}

	// 3) Overwrite the same key
	if err := store.Put("foo", "baz"); err != nil {
		t.Fatalf("unexpected error on Put overwrite: %v", err)
	}
	got, err = store.Get("foo")
	if err != nil {
		t.Fatalf("unexpected error on Get after overwrite: %v", err)
	}
	if got != "baz" {
		t.Fatalf("Get after overwrite returned %q, want %q", got, "baz")
	}

	// 4) Delete the key, then Get should error with ErrMsgKeyNotFound
	if err := store.Delete("foo"); err != nil {
		t.Fatalf("unexpected error on Delete: %v", err)
	}
	_, err = store.Get("foo")
	if err == nil {
		t.Fatal("expected error after Delete, got nil")
	}
	// optionally check error message
	if !errors.Is(err, errors.New(ErrMsgKeyNotFound)) && err.Error() != ErrMsgKeyNotFound {
		t.Fatalf("expected ErrMsgKeyNotFound, got %v", err)
	}
}

func newTestServer() *Server {
	logger := log.New(io.Discard, "", 0)
	return NewServer(NewInMemoryKVS(), logger)
}

func TestServerIntegration(t *testing.T) {
	// Create server backed by a fresh in-memory Store
	srv := newTestServer()

	// 1) GET missing key → 404
	{
		req := httptest.NewRequest(http.MethodGet, "/kv/foo", nil)
		rec := httptest.NewRecorder()
		srv.ServeHTTP(rec, req)

		if rec.Code != http.StatusNotFound {
			t.Errorf("GET missing: expected 404, got %d", rec.Code)
		}
	}

	// 2) POST new key → 202 Accepted
	{
		req := httptest.NewRequest(http.MethodPost, "/kv/foo", strings.NewReader("bar"))
		rec := httptest.NewRecorder()
		srv.ServeHTTP(rec, req)

		if rec.Code != http.StatusAccepted {
			t.Errorf("POST create: expected 202, got %d", rec.Code)
		}
	}

	// 3) GET existing key → 200 OK + "bar"
	{
		req := httptest.NewRequest(http.MethodGet, "/kv/foo", nil)
		rec := httptest.NewRecorder()
		srv.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("GET existing: expected 200, got %d", rec.Code)
		}
		got, _ := io.ReadAll(rec.Body)
		if gotStr := string(got); gotStr != "bar" {
			t.Errorf("GET existing: expected body %q, got %q", "bar", gotStr)
		}
	}

	// 4) DELETE existing key → 200 OK
	{
		req := httptest.NewRequest(http.MethodDelete, "/kv/foo", nil)
		rec := httptest.NewRecorder()
		srv.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("DELETE existing: expected 200, got %d", rec.Code)
		}
	}

	// 5) GET after delete → 404 Not Found
	{
		req := httptest.NewRequest(http.MethodGet, "/kv/foo", nil)
		rec := httptest.NewRecorder()
		srv.ServeHTTP(rec, req)

		if rec.Code != http.StatusNotFound {
			t.Errorf("GET after delete: expected 404, got %d", rec.Code)
		}
	}

	// 6) DELETE missing key → 404 Not Found
	{
		req := httptest.NewRequest(http.MethodDelete, "/kv/nonexistent", nil)
		rec := httptest.NewRecorder()
		srv.ServeHTTP(rec, req)

		if rec.Code != http.StatusNotFound {
			t.Errorf("DELETE missing: expected 404, got %d", rec.Code)
		}
	}
}

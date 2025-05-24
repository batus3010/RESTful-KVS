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
	database, cleanDatabase := createTempFileSystem(t, "")
	defer cleanDatabase()
	store := NewFileSystemKVStore(database)

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
		t.Fatal("expected error ErrMsgKeyNotFound after Delete, got nil")
	}
	if !errors.Is(err, errors.New(ErrMsgKeyNotFound)) && err.Error() != ErrMsgKeyNotFound {
		t.Fatalf("expected ErrMsgKeyNotFound, got %v", err)
	}
}

// --- HTTP server integration ---

// newTestServer gives you a Server with an empty InMemoryKVS and a silent logger.
func newTestServer() *Server {
	logger := log.New(io.Discard, "", 0)
	return NewServer(NewInMemoryKVS(), logger)
}

func TestHTTPIntegration(t *testing.T) {
	srv := newTestServer()

	// 1) GET missing key → 404
	{
		req := httptest.NewRequest(http.MethodGet, "/kv/foo", nil)
		rec := httptest.NewRecorder()
		srv.ServeHTTP(rec, req)
		assertStatus(t, rec.Code, http.StatusNotFound)
	}

	// 2) POST new key → 202 Accepted
	{
		req := httptest.NewRequest(http.MethodPost, "/kv/foo", strings.NewReader("bar"))
		rec := httptest.NewRecorder()
		srv.ServeHTTP(rec, req)
		assertStatus(t, rec.Code, http.StatusAccepted)
	}

	// 3) GET existing key → 200 OK + "bar"
	{
		req := httptest.NewRequest(http.MethodGet, "/kv/foo", nil)
		rec := httptest.NewRecorder()
		srv.ServeHTTP(rec, req)
		assertStatus(t, rec.Code, http.StatusOK)
		assertBody(t, rec.Body, "bar")
	}

	// 4) DELETE existing key → 200 OK
	{
		req := httptest.NewRequest(http.MethodDelete, "/kv/foo", nil)
		rec := httptest.NewRecorder()
		srv.ServeHTTP(rec, req)
		assertStatus(t, rec.Code, http.StatusOK)
	}

	// 5) GET after delete → 404 Not Found
	{
		req := httptest.NewRequest(http.MethodGet, "/kv/foo", nil)
		rec := httptest.NewRecorder()
		srv.ServeHTTP(rec, req)
		assertStatus(t, rec.Code, http.StatusNotFound)
	}

	// 6) DELETE missing key → 404 Not Found
	{
		req := httptest.NewRequest(http.MethodDelete, "/kv/nonexistent", nil)
		rec := httptest.NewRecorder()
		srv.ServeHTTP(rec, req)
		assertStatus(t, rec.Code, http.StatusNotFound)
	}
}

func TestListAllKeys(t *testing.T) {
	// seed store with two entries
	store := NewInMemoryKVS()
	store.Put("foo", "bar")
	store.Put("baz", "qux")

	srv := newTestServer()
	// replace the store inside srv so it's pre-populated
	srv.Store = store

	// GET /all → JSON array of KVPair
	req := httptest.NewRequest(http.MethodGet, "/all", nil)
	rec := httptest.NewRecorder()
	srv.ServeHTTP(rec, req)

	// should return 200 + application/json
	assertStatus(t, rec.Code, http.StatusOK)
	assertContentType(t, rec, jsonContentType)

	// decode and compare
	got := getTableFromResponse(t, rec.Body)

	expected := Table{
		{Key: "foo", Value: "bar"},
		{Key: "baz", Value: "qux"},
	}
	assertTable(t, got, expected)
}

func assertBody(t testing.TB, body io.Reader, want string) {
	t.Helper()
	b, err := io.ReadAll(body)
	if err != nil {
		t.Fatalf("reading body: %v", err)
	}
	if got := string(b); got != want {
		t.Errorf("body: got %q, want %q", got, want)
	}
}

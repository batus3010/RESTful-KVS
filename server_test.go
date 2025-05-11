package kvs

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

type StubStore struct {
	kv map[string]string
}

func (store StubStore) Get(key string) (string, error) {
	value := store.kv[key]
	return value, nil
}

func (store StubStore) Put(key string, value string) error {
	store.kv[key] = value
	return nil
}

func TestGet(t *testing.T) {
	t.Run("Get returns 404 on missing key", func(t *testing.T) {
		store := &StubStore{}
		server := NewServer(store)
		request := httptest.NewRequest(http.MethodGet, "/kv/foo", nil)
		response := httptest.NewRecorder()
		server.ServeHTTP(response, request)
		assertStatus(t, response.Code, http.StatusNotFound)
	})
}

func assertStatus(t testing.TB, got, want int) {
	t.Helper()
	if got != want {
		t.Errorf("did not get correct status, got %d, want %d", got, want)
	}
}

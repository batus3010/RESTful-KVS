package kvs

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type StubStore struct {
	kv map[string]string
}

func (store *StubStore) Get(key string) (string, error) {
	value := store.kv[key]
	return value, nil
}

func (store *StubStore) Put(key string, value string) error {
	store.kv[key] = value
	return nil
}

func TestGet(t *testing.T) {
	t.Run("Get returns 404 on missing key", func(t *testing.T) {
		server, response := newTestServer(map[string]string{})
		request := httptest.NewRequest(http.MethodGet, "/kv/foo", nil)
		server.ServeHTTP(response, request)
		assertStatus(t, response.Code, http.StatusNotFound)
	})
	t.Run("Get on existing key return 200 and value", func(t *testing.T) {
		server, response := newTestServer(map[string]string{"foo": "bar"})
		request := httptest.NewRequest(http.MethodGet, "/kv/foo", nil)
		server.ServeHTTP(response, request)
		assertStatus(t, response.Code, http.StatusOK)
		assertEqual(t, response.Body.String(), "bar")
	})
}

func TestPut(t *testing.T) {
	t.Run("Put returns 201", func(t *testing.T) {
		server, response := newTestServer(map[string]string{})
		request := httptest.NewRequest(http.MethodPost, "/kv/foo", strings.NewReader("bar"))
		server.ServeHTTP(response, request)
		assertStatus(t, response.Code, http.StatusAccepted)
	})
	t.Run("Put new value then Get return that value", func(t *testing.T) {
		server, response := newTestServer(map[string]string{})
		request := httptest.NewRequest(http.MethodPost, "/kv/foo", strings.NewReader("bar"))
		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusAccepted)

		getRequest := httptest.NewRequest(http.MethodGet, "/kv/foo", nil)
		getResponse := httptest.NewRecorder()
		server.ServeHTTP(getResponse, getRequest)
		assertStatus(t, getResponse.Code, http.StatusOK)
		assertEqual(t, getResponse.Body.String(), "bar")
	})
}

func newTestServer(initial map[string]string) (*Server, *httptest.ResponseRecorder) {
	store := &StubStore{kv: initial}
	srv := NewServer(store)
	response := httptest.NewRecorder()
	return srv, response
}

func assertStatus(t testing.TB, got, want int) {
	t.Helper()
	if got != want {
		t.Errorf("did not get correct status, got %d, want %d", got, want)
	}
}

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

type StubStore struct {
	kv  map[string]string
	all []KVPair
}

func (store *StubStore) Get(key string) (string, error) {
	if val, ok := store.kv[key]; ok {
		return val, nil
	}
	return "", errors.New("value not found")
}

func (store *StubStore) Put(key string, value string) error {
	store.kv[key] = value
	return nil
}

func (store *StubStore) Delete(key string) error {
	delete(store.kv, key)
	return nil
}

func TestGet(t *testing.T) {
	t.Run("Get returns 404 on missing key", func(t *testing.T) {
		server, response := newTestServerWithStubStore(StubStore{})
		request := httptest.NewRequest(http.MethodGet, "/kv/foo", nil)
		server.ServeHTTP(response, request)
		assertStatus(t, response.Code, http.StatusNotFound)
	})
	t.Run("Get on existing key return 200 and value", func(t *testing.T) {
		store := StubStore{kv: map[string]string{
			"foo": "bar",
		}}
		server, response := newTestServerWithStubStore(store)
		request := httptest.NewRequest(http.MethodGet, "/kv/foo", nil)
		server.ServeHTTP(response, request)
		assertStatus(t, response.Code, http.StatusOK)
		assertEqual(t, response.Body.String(), "bar")
	})
}

func TestPut(t *testing.T) {
	t.Run("Put returns 201", func(t *testing.T) {
		server, response := newTestServerWithStubStore(StubStore{kv: map[string]string{}})
		request := httptest.NewRequest(http.MethodPost, "/kv/foo", strings.NewReader("bar"))
		server.ServeHTTP(response, request)
		assertStatus(t, response.Code, http.StatusAccepted)
	})
	t.Run("Put new value then Get return that value", func(t *testing.T) {
		server, response := newTestServerWithStubStore(StubStore{kv: map[string]string{}})
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

func TestDelete(t *testing.T) {
	t.Run("Delete non-existing key returns 404", func(t *testing.T) {
		server, response := newTestServerWithStubStore(StubStore{})
		request := httptest.NewRequest(http.MethodDelete, "/kv/foo", nil)
		server.ServeHTTP(response, request)
		assertStatus(t, response.Code, http.StatusNotFound)
	})
	t.Run("Delete existing key", func(t *testing.T) {
		store := StubStore{kv: map[string]string{
			"foo": "bar",
		}}
		server, response := newTestServerWithStubStore(store)
		request := httptest.NewRequest(http.MethodDelete, "/kv/foo", nil)
		server.ServeHTTP(response, request)
		assertStatus(t, response.Code, http.StatusOK)
		getRequest := httptest.NewRequest(http.MethodGet, "/kv/foo", nil)
		getResponse := httptest.NewRecorder()
		server.ServeHTTP(getResponse, getRequest)
		assertStatus(t, getResponse.Code, http.StatusNotFound)
	})
}

func TestRestKVS(t *testing.T) {
	t.Run("Method not allowed should returns StatusMethodNotAllowed", func(t *testing.T) {
		store := StubStore{kv: map[string]string{}}
		server, response := newTestServerWithStubStore(store)
		request := httptest.NewRequest(http.MethodPut, "/kv/foo", nil)
		server.ServeHTTP(response, request)
		assertStatus(t, response.Code, http.StatusMethodNotAllowed)
	})
}

//func TestAll(t *testing.T) {
//	t.Run("return key value table as JSON", func(t *testing.T) {
//		wantedTable := []KVPair{
//			{Key: "foo", Value: "bar"},
//			{Key: "bar", Value: "baz"},
//		}
//		store := StubStore{nil, wantedTable}
//		server :=
//	})
//}

func newTestServerWithStubStore(store StubStore) (*Server, *httptest.ResponseRecorder) {
	silentLog := log.New(io.Discard, "", 0)
	srv := NewServer(&store, silentLog)
	response := httptest.NewRecorder()
	return srv, response
}

func assertStatus(t testing.TB, got, want int) {
	t.Helper()
	if got != want {
		t.Errorf("did not get correct status, got %d, want %d", got, want)
	}
}

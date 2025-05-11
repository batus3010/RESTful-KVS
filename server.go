package kvs

import (
	"io"
	"net/http"
	"strings"
)

type Server struct {
	store KeyValueStore
	http.Handler
}

func NewServer(store KeyValueStore) *Server {
	p := new(Server)
	p.store = store
	router := http.NewServeMux()
	router.Handle("/kv/", http.HandlerFunc(p.storeHandler))
	p.Handler = router
	return p
}

func (p *Server) storeHandler(w http.ResponseWriter, r *http.Request) {
	key := strings.TrimPrefix(r.URL.Path, "/kv/")
	switch r.Method {
	case http.MethodGet:
		p.handleGet(w, key)
	case http.MethodPost:
		p.handlePost(w, r, key)
	case http.MethodDelete:
		p.handleDelete(w, key)
	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

func (p *Server) handleGet(w http.ResponseWriter, key string) {
	value, _ := p.store.Get(key)
	if value == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.Write([]byte(value))
}

func (p *Server) handlePost(w http.ResponseWriter, r *http.Request, key string) {
	body, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		return
	}
	value := string(body)
	if err := p.store.Put(key, value); err != nil {
		// if your Put ever returns an error, report it
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)
}

func (p *Server) handleDelete(w http.ResponseWriter, key string) {
	_, err := p.store.Get(key)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	err = p.store.Delete(key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
}

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

func (p *Server) storeHandler(w http.ResponseWriter, req *http.Request) {
	key := strings.TrimPrefix(req.URL.Path, "/kv/")
	switch req.Method {
	case http.MethodGet:
		p.showValue(w, key)
	case http.MethodPost:
		body, err := io.ReadAll(req.Body)
		if err != nil {
			http.Error(w, "Error reading request body", http.StatusBadRequest)
			return
		}
		defer req.Body.Close()
		value := string(body)
		p.store.Put(key, value)
		w.WriteHeader(http.StatusAccepted)
	}
}

func (p *Server) showValue(w http.ResponseWriter, key string) {
	value, _ := p.store.Get(key)
	if value == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.Write([]byte(value))
}

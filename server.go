package kvs

import (
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
		value, _ := p.store.Get(key)
		if value == "" {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.Write([]byte(value))
	}
}

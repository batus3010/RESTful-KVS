package kvs

import (
	"net/http"
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
	switch req.Method {
	case http.MethodGet:
		http.NotFound(w, req)
	}
}

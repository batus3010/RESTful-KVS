package kvs

import "net/http"

type Server struct {
	store KeyValueStore
	http.Handler
}

func NewServer(store KeyValueStore) *Server {
	return &Server{store: store, Handler: http.NewServeMux()}
}

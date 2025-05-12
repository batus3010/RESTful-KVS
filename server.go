package kvs

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

type Server struct {
	Store  KeyValueStore
	logger *log.Logger
	http.Handler
}

type KVPair struct {
	Key   string
	Value string
}

func NewServer(store KeyValueStore, logger *log.Logger) *Server {
	srv := &Server{
		Store:  store,
		logger: logger,
	}

	mux := http.NewServeMux()
	// wrap storeHandler with requestLogger
	mux.Handle("/kv/", srv.requestLogger(http.HandlerFunc(srv.storeHandler)))
	mux.Handle("/all", srv.requestLogger(http.HandlerFunc(srv.allHandler)))
	srv.Handler = mux
	return srv
}

func (srv *Server) storeHandler(w http.ResponseWriter, r *http.Request) {
	key := strings.TrimPrefix(r.URL.Path, "/kv/")
	switch r.Method {
	case http.MethodGet:
		srv.handleGet(w, key)
	case http.MethodPost:
		srv.handlePost(w, r, key)
	case http.MethodDelete:
		srv.handleDelete(w, key)
	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

func (srv *Server) allHandler(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(srv.getAllKeys())
	w.WriteHeader(http.StatusOK)
}

func (srv *Server) getAllKeys() []KVPair {
	return []KVPair{
		{"key1", "value1"},
	}
}

func (srv *Server) handleGet(w http.ResponseWriter, key string) {
	value, _ := srv.Store.Get(key)
	if value == "" {
		srv.logger.Printf("GET %q → not found", key)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	srv.logger.Printf("GET %q → %q", key, value)
	w.Write([]byte(value))
}

func (srv *Server) handlePost(w http.ResponseWriter, r *http.Request, key string) {
	body, err := io.ReadAll(r.Body)
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			srv.logger.Printf("POST %q → error closing request body: %v", key, err)
		}
	}(r.Body)
	if err != nil {
		srv.logger.Printf("POST %q → read error: %v", key, err)
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		return
	}
	value := string(body)
	if err := srv.Store.Put(key, value); err != nil {
		// if your Put ever returns an error, report it
		srv.logger.Printf("POST %q → store error: %v", key, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	srv.logger.Printf("POST %q → set to %q", key, value)
	w.WriteHeader(http.StatusAccepted)
}

func (srv *Server) handleDelete(w http.ResponseWriter, key string) {
	_, err := srv.Store.Get(key)
	if err != nil {
		srv.logger.Printf("DELETE %q → not found", key)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	err = srv.Store.Delete(key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	srv.logger.Printf("DELETE %q → deleted", key)
	w.WriteHeader(http.StatusOK)
}

func (srv *Server) requestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// wrap the ResponseWriter to capture status code
		lrw := &loggingResponseWriter{w, http.StatusOK}
		start := time.Now()

		next.ServeHTTP(lrw, r)

		srv.logger.Printf(
			"%s %s %s → %d (%s)",
			r.RemoteAddr,
			r.Method,
			r.URL.Path,
			lrw.statusCode,
			time.Since(start),
		)
	})
}

// helper to capture WriteHeader calls
type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

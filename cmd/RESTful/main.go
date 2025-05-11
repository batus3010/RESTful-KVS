package main

import (
	kvs "KVS"
	"log"
	"net/http"
	"os"
)

func main() {
	// log to stdout with timestamp
	logger := log.New(os.Stdout, "[KVS] ", log.LstdFlags)

	store := kvs.NewInMemoryKVS()
	server := kvs.NewServer(store, logger)

	logger.Printf("starting server on :5000")
	log.Fatal(http.ListenAndServe(":5000", server))
}

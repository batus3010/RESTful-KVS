package main

import (
	kvs "KVS"
	"log"
	"net/http"
)

func main() {
	server := kvs.NewServer(kvs.NewInMemoryKVS()) // <— this wires up Handler
	log.Fatal(http.ListenAndServe(":6969", server))
}

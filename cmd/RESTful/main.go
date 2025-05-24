package main

import (
	kvs "KVS"
	"log"
	"net/http"
	"os"
)

const dbFileName = "KVS.db.json"

func main() {
	// log to stdout with timestamp
	logger := log.New(os.Stdout, "[KVS] ", log.LstdFlags)

	//store := kvs.NewInMemoryKVS()
	//server := kvs.NewServer(store, logger)
	//
	//logger.Printf("starting server on :5000")
	//log.Fatal(http.ListenAndServe(":5000", server))
	dbFile, err := os.OpenFile(dbFileName,
		os.O_RDWR|os.O_CREATE,
		0o600, // owner read/write only
	)
	if err != nil {
		logger.Fatalf("failed to open database file: %v", err)
	}
	defer dbFile.Close()
	store, err := kvs.NewFileSystemKVStore(dbFile)

	if err != nil {
		logger.Fatalf("failed to initialize store: %v", err)
	}

	server := kvs.NewServer(store, logger)
	logger.Printf("starting server on :5000")
	if err := http.ListenAndServe(":5000", server); err != nil {
		log.Fatalf("could not listen on port 5000 %v", err)
	}
}

package kvs

import (
	"encoding/json"
	"io"
)

type FileKVStore struct {
	database io.Reader
}

func (f *FileKVStore) GetTable() []KVPair {
	var table []KVPair
	json.NewDecoder(f.database).Decode(&table)
	return table
}

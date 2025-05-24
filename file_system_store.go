package kvs

import (
	"io"
)

type FileKVStore struct {
	database io.Reader
}

func (f *FileKVStore) GetTable() []KVPair {
	table, _ := NewTable(f.database)
	return table
}

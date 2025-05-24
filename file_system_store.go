package kvs

import (
	"io"
)

type FileSystemKVStore struct {
	database io.ReadSeeker
}

func (f *FileSystemKVStore) GetTable() []KVPair {
	f.database.Seek(0, io.SeekStart)
	table, _ := NewTable(f.database)
	return table
}

func (f *FileSystemKVStore) GetValueOf(key string) string {
	var value string
	for _, k := range f.GetTable() {
		if k.Key == key {
			value = k.Value
			break
		}
	}
	return value
}

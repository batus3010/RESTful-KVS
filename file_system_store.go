package kvs

import (
	"encoding/json"
	"io"
)

type FileSystemKVStore struct {
	database io.ReadWriteSeeker
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

func (f *FileSystemKVStore) Update(key string, value string) {
	table := f.GetTable()
	for i, k := range table {
		if k.Key == key {
			table[i].Value = value
		}
	}
	f.database.Seek(0, io.SeekStart)
	json.NewEncoder(f.database).Encode(table)
}

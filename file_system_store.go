package kvs

import (
	"encoding/json"
	"io"
)

type FileSystemKVStore struct {
	database io.ReadWriteSeeker
}

func (f *FileSystemKVStore) GetTable() Table {
	f.database.Seek(0, io.SeekStart)
	table, _ := NewTable(f.database)
	return table
}

func (f *FileSystemKVStore) GetValueOf(key string) string {
	pair := f.GetTable().Find(key)
	if pair != nil {
		return pair.Value
	}
	return ""
}

func (f *FileSystemKVStore) Update(key string, value string) {
	table := f.GetTable()
	pair := table.Find(key)
	if pair != nil {
		pair.Value = value
	}
	f.database.Seek(0, io.SeekStart)
	json.NewEncoder(f.database).Encode(table)
}

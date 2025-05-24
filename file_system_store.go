package kvs

import (
	"encoding/json"
	"errors"
	"io"
	"os"
)

type FileSystemKVStore struct {
	database io.Writer
	table    Table
}

func NewFileSystemKVStore(file *os.File) *FileSystemKVStore {
	file.Seek(0, io.SeekStart)
	table, _ := NewTable(file)

	writer := &rewindableWriter{file: file}
	return &FileSystemKVStore{
		database: writer,
		table:    table,
	}
}

func (f *FileSystemKVStore) GetTable() Table {
	return f.table
}

func (f *FileSystemKVStore) Get(key string) (string, error) {
	pair := f.table.Find(key)
	if pair != nil {
		return pair.Value, nil
	}
	return "", errors.New(ErrMsgKeyNotFound)
}

func (f *FileSystemKVStore) Put(key string, value string) error {
	pair := f.table.Find(key)
	if pair != nil {
		pair.Value = value
	} else {
		f.table = append(f.table, KVPair{key, value})
	}
	return json.NewEncoder(f.database).Encode(f.table)
}

func (f *FileSystemKVStore) Delete(key string) error {
	if !f.table.Remove(key) {
		return errors.New(ErrMsgKeyNotFound)
	}
	// one call to rewrite the file
	return json.NewEncoder(f.database).Encode(f.table)
}

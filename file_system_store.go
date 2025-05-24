package kvs

import (
	"encoding/json"
	"errors"
	"io"
	"os"
)

type FileSystemKVStore struct {
	database *json.Encoder
	table    Table
}

func NewFileSystemKVStore(file *os.File) *FileSystemKVStore {
	file.Seek(0, io.SeekStart)
	table, _ := NewTable(file)

	return &FileSystemKVStore{
		database: json.NewEncoder(&rewindableWriter{file}),
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
	err := f.database.Encode(f.table)
	if err != nil {
		return err
	}
	return nil
}

func (f *FileSystemKVStore) Delete(key string) error {
	if !f.table.Remove(key) {
		return errors.New(ErrMsgKeyNotFound)
	}
	// one call to rewrite the file
	err := f.database.Encode(f.table)
	if err != nil {
		return err
	}
	return nil
}

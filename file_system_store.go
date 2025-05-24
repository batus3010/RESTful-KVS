package kvs

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
)

type FileSystemKVStore struct {
	database *json.Encoder
	table    Table
}

func NewFileSystemKVStore(file *os.File) (*FileSystemKVStore, error) {
	file.Seek(0, io.SeekStart)

	info, err := file.Stat()

	if err != nil {
		return nil, fmt.Errorf("problem getting file info from file %s, %v", file.Name(), err)
	}

	if info.Size() == 0 {
		file.Write([]byte("[]"))
		file.Seek(0, io.SeekStart)
	}

	table, err := NewTable(file)

	if err != nil {
		return nil, fmt.Errorf("problem loading player store from file %s, %v", file.Name(), err)
	}

	return &FileSystemKVStore{
		database: json.NewEncoder(&rewindableWriter{file}),
		table:    table,
	}, nil
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

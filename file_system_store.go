package kvs

import (
	"encoding/json"
	"errors"
	"io"
)

type FileSystemKVStore struct {
	database io.ReadWriteSeeker
	table    Table
}

func NewFileSystemKVStore(database io.ReadWriteSeeker) *FileSystemKVStore {
	database.Seek(0, io.SeekStart)
	table, _ := NewTable(database)
	return &FileSystemKVStore{
		database: database,
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
	f.database.Seek(0, io.SeekStart)
	json.NewEncoder(f.database).Encode(f.table)
	return nil
}

func (f *FileSystemKVStore) Delete(key string) error {
	if !f.table.Remove(key) {
		return errors.New(ErrMsgKeyNotFound)
	}

	// Seek back to the beginning of file
	if _, err := f.database.Seek(0, io.SeekStart); err != nil {
		return err
	}

	// Truncate the file so old JSON is wiped
	if t, ok := f.database.(interface{ Truncate(int64) error }); ok {
		if err := t.Truncate(0); err != nil {
			return err
		}
		// seek again just in case
		if _, err := f.database.Seek(0, io.SeekStart); err != nil {
			return err
		}
	}

	// Encode the updated table back to disk
	return json.NewEncoder(f.database).Encode(f.table)
}

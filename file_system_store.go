package kvs

import (
	"encoding/json"
	"errors"
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

func (f *FileSystemKVStore) Get(key string) (string, error) {
	pair := f.GetTable().Find(key)
	if pair != nil {
		return pair.Value, nil
	}
	return "", errors.New(ErrMsgKeyNotFound)
}

func (f *FileSystemKVStore) Put(key string, value string) error {
	table := f.GetTable()
	pair := table.Find(key)
	if pair != nil {
		pair.Value = value
	} else {
		table = append(table, KVPair{key, value})
	}
	f.database.Seek(0, io.SeekStart)
	json.NewEncoder(f.database).Encode(table)
	return nil
}

func (f *FileSystemKVStore) Delete(key string) error {
	// 1) Load current entries
	table := f.GetTable()

	// 2) Build a new table without the deleted key
	var newTable Table
	var found bool
	for _, pair := range table {
		if pair.Key == key {
			found = true
			continue
		}
		newTable = append(newTable, pair)
	}

	// 3) Key not found → error
	if !found {
		return errors.New(ErrMsgKeyNotFound)
	}

	// 4) Overwrite the file from the beginning
	// Seek to start
	if _, err := f.database.Seek(0, io.SeekStart); err != nil {
		return err
	}
	// Truncate if supported (e.g. os.File)
	if t, ok := f.database.(interface{ Truncate(int64) error }); ok {
		if err := t.Truncate(0); err != nil {
			return err
		}
		// Seek again just in case
		if _, err := f.database.Seek(0, io.SeekStart); err != nil {
			return err
		}
	}

	// Encode the new table
	return json.NewEncoder(f.database).Encode(newTable)
}

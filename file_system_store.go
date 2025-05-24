package kvs

import (
	"encoding/json"
	"errors"
	"io"
)

type FileSystemKVStore struct {
	Database io.ReadWriteSeeker
}

func (f *FileSystemKVStore) GetTable() Table {
	f.Database.Seek(0, io.SeekStart)
	table, _ := NewTable(f.Database)
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
	f.Database.Seek(0, io.SeekStart)
	json.NewEncoder(f.Database).Encode(table)
	return nil
}

func (f *FileSystemKVStore) Delete(key string) error {
	table := f.GetTable()
	if !table.Remove(key) {
		return errors.New(ErrMsgKeyNotFound)
	}

	// Seek back to the beginning of file
	if _, err := f.Database.Seek(0, io.SeekStart); err != nil {
		return err
	}

	// Truncate the file so old JSON is wiped
	if t, ok := f.Database.(interface{ Truncate(int64) error }); ok {
		if err := t.Truncate(0); err != nil {
			return err
		}
		// seek again just in case
		if _, err := f.Database.Seek(0, io.SeekStart); err != nil {
			return err
		}
	}

	// Encode the updated table back to disk
	return json.NewEncoder(f.Database).Encode(table)
}

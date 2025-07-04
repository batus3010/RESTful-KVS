package kvs

import (
	"os"
	"testing"
)

func TestFileSystemStore(t *testing.T) {
	t.Run("works with an empty file", func(t *testing.T) {
		database, cleanDatabase := createTempFileSystem(t, "")
		defer cleanDatabase()

		_, err := NewFileSystemKVStore(database)

		assertNoError(t, err)
	})

	t.Run("KV table from a reader", func(t *testing.T) {
		database, cleanDatabase := createTempFileSystem(t, `[
			{"Key": "key1", "Value": "value1"},
			{"Key": "key2", "Value": "value2"}]`)
		defer cleanDatabase()
		store, err := NewFileSystemKVStore(database)
		assertNoError(t, err)

		got := store.GetTable()
		want := []KVPair{
			{"key1", "value1"},
			{"key2", "value2"},
		}
		assertTable(t, got, want)

		// read again to test for multiple reads
		got = store.GetTable()
		assertTable(t, got, want)
	})

	t.Run("get value from a reader", func(t *testing.T) {
		database, cleanDatabase := createTempFileSystem(t, `[
			{"Key": "key1", "Value": "value1"},
			{"Key": "key2", "Value": "value2"}]`)
		defer cleanDatabase()
		store, err := NewFileSystemKVStore(database)
		assertNoError(t, err)

		got, _ := store.Get("key1")
		want := "value1"
		assertEqual(t, got, want)
	})

	t.Run("update value for existing key", func(t *testing.T) {
		database, cleanDatabase := createTempFileSystem(t, `[
			{"Key": "key1", "Value": "old value"},
			{"Key": "key2", "Value": "value2"}]`)
		defer cleanDatabase()

		store, err := NewFileSystemKVStore(database)
		assertNoError(t, err)

		store.Put("key1", "new value")
		got, _ := store.Get("key1")
		want := "new value"
		assertEqual(t, got, want)
	})

	t.Run("store value for new key", func(t *testing.T) {
		database, cleanDatabase := createTempFileSystem(t, `[
			{"Key": "key1", "Value": "old value"},
			{"Key": "key2", "Value": "value2"}]`)
		defer cleanDatabase()

		store, err := NewFileSystemKVStore(database)
		assertNoError(t, err)

		store.Put("key3", "value3")
		got, _ := store.Get("key3")
		want := "value3"
		assertEqual(t, got, want)
	})
}

func createTempFileSystem(t testing.TB, initialData string) (*os.File, func()) {
	t.Helper()
	tmpfile, err := os.CreateTemp("", "db")
	if err != nil {
		t.Fatalf("could not create temp file %v", err)
	}
	tmpfile.Write([]byte(initialData))
	removeFile := func() {
		tmpfile.Close()
		os.Remove(tmpfile.Name())
	}

	return tmpfile, removeFile
}

package kvs

import (
	"strings"
	"testing"
)

func TestFileSystemStore(t *testing.T) {
	t.Run("KV table from a reader", func(t *testing.T) {
		database := strings.NewReader(`[
			{"Key": "key1", "Value": "value1"},
			{"Key": "key2", "Value": "value2"}]`)
		store := FileSystemKVStore{database}

		got := store.GetTable()
		want := []KVPair{
			{"key1", "value1"},
			{"key2", "value2"},
		}
		assertTable(t, got, want)

		got = store.GetTable()
		assertTable(t, got, want)
	})

	t.Run("get value from a reader", func(t *testing.T) {
		database := strings.NewReader(`[
			{"Key": "key1", "Value": "value1"},
			{"Key": "key2", "Value": "value2"}]`)
		store := FileSystemKVStore{database}
		got := store.GetValueOf("key1")
		want := "value1"
		assertEqual(t, got, want)
	})
}

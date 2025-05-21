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
		store := FileKVStore{database}

		got := store.GetTable()
		want := []KVPair{
			{"key1", "value1"},
			{"key2", "value2"},
		}
		assertTable(t, got, want)
	})
}

package Key_Value_Store

import "testing"

func TestKVS(t *testing.T) {
	t.Run("GET non-exist key returns empty string and nil", func(t *testing.T) {
		store := NewKeyValueStore()
		got, _ := store.Get("foo")
		want := ""
		if got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})
}

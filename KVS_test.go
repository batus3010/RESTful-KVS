package Key_Value_Store

import "testing"

func TestKVS(t *testing.T) {
	t.Run("Get non-exist key returns empty string and error", func(t *testing.T) {
		kvs := NewInMemoryKVS()
		got, err := kvs.Get("foo")
		want := ""
		assertEqual(t, got, want)
		assertError(t, err)
	})
	t.Run("Put value and then Get returns value and nil", func(t *testing.T) {})
}

func assertError(t *testing.T, err error) {
	t.Helper()
	if err == nil {
		t.Fatalf("Expected to get error but get %v instead", err)
	}
}

func assertEqual[T comparable](t *testing.T, got, want T) {
	t.Helper()
	if got != want {
		t.Errorf("got %v want %v", got, want)
	}
}

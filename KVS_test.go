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
	t.Run("Put value and then Get returns value and nil", func(t *testing.T) {
		kvs := NewInMemoryKVS()
		err := kvs.Put("foo", "bar")
		assertNoError(t, err)
		got, err := kvs.Get("foo")
		want := "bar"
		assertEqual(t, got, want)
	})
}

func assertError(t *testing.T, err error) {
	t.Helper()
	if err == nil {
		t.Fatalf("Expected to get error but get %v instead", err)
	}
}

func assertNoError(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("Expected no error but got %v instead", err)
	}
}

func assertEqual[T comparable](t *testing.T, got, want T) {
	t.Helper()
	if got != want {
		t.Errorf("got '%v' want '%v'", got, want)
	}
}

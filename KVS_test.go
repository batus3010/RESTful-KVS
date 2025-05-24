package kvs

import "testing"

func TestKVS(t *testing.T) {
	database, cleanDatabase := createTempFileSystem(t, "")
	defer cleanDatabase()
	cases := []struct {
		name      string
		setup     func(kvs KeyValueStore)
		key       string
		wantValue string
		wantErr   bool
	}{
		{
			name:    "missing key",
			setup:   func(kvs KeyValueStore) {},
			key:     "foo",
			wantErr: true,
		},
		{
			name:      "single put",
			setup:     func(kvs KeyValueStore) { kvs.Put("foo", "bar") },
			key:       "foo",
			wantValue: "bar",
			wantErr:   false,
		},
		{
			name: "overwrite key",
			setup: func(kvs KeyValueStore) {
				kvs.Put("foo", "bar")
				kvs.Put("foo", "baz")
			},
			key:       "foo",
			wantValue: "baz",
			wantErr:   false,
		},
		// Add more cases ...
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			kvs, err := NewFileSystemKVStore(database)
			assertNoError(t, err)
			tc.setup(kvs)
			got, err := kvs.Get(tc.key)
			if tc.wantErr {
				assertError(t, err)
			} else {
				assertNoError(t, err)
				assertEqual(t, got, tc.wantValue)
			}
		})
	}
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

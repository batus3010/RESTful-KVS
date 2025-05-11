package Key_Value_Store

import "errors"

const (
	ErrMsgKeyNotFound = "key not found"
)

type KeyValueStore interface {
	Put(key, value string) error
	Get(key string) (string, error)
}

type InMemoryKVS struct {
	Store map[string]string
}

func NewInMemoryKVS() KeyValueStore {
	return &InMemoryKVS{
		Store: make(map[string]string),
	}
}

func (kvs *InMemoryKVS) Put(key, value string) error {
	kvs.Store[key] = value
	return nil
}

func (kvs *InMemoryKVS) Get(key string) (string, error) {
	if val, ok := kvs.Store[key]; ok {
		return val, nil
	}
	return "", errors.New(ErrMsgKeyNotFound)
}

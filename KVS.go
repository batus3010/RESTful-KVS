package kvs

import "errors"

const (
	ErrMsgKeyNotFound = "key not found"
)

type KeyValueStore interface {
	Put(key, value string) error
	Get(key string) (string, error)
	Delete(key string) error
	GetTable() []KVPair
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

func (kvs *InMemoryKVS) Delete(key string) error {
	delete(kvs.Store, key)
	return nil
}

func (kvs *InMemoryKVS) GetTable() []KVPair {
	var KVPairs []KVPair
	for key, value := range kvs.Store {
		KVPairs = append(KVPairs, KVPair{key, value})
	}
	return KVPairs
}

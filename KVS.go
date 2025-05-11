package Key_Value_Store

type KeyValueStore interface {
	Put(key, value string) error
	Get(key string) (string, error)
}

type InMemoryKVS struct {
	store KeyValueStore
}

func NewInMemoryKVS() KeyValueStore {
	return &InMemoryKVS{}
}

func (kvs *InMemoryKVS) Put(key, value string) error {
	return nil
}

func (kvs *InMemoryKVS) Get(key string) (string, error) {
	return "", nil
}

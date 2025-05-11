package Key_Value_Store

type KeyValueStore interface {
	Put(key, value string) error
	Get(key string) (string, error)
}

type KVS struct{}

func NewKeyValueStore() KeyValueStore {
	return &KVS{}
}

func (kvs *KVS) Put(key, value string) error {
	return nil
}

func (kvs *KVS) Get(key string) (string, error) {
	return "", nil
}

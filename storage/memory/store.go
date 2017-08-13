package memory

import "sync"

// Provides an in memory, non persistent, implementation of Store interface.

type Store struct {
	sync.RWMutex
	m map[string]string
}

func NewStore() *Store {
	return &Store{m: make(map[string]string)}
}

func (store *Store) Get(key string) ([]byte, error) {
	store.RLock()
	value := store.m[key]
	store.RUnlock()

	return []byte(value), nil
}

func (store *Store) Put(key string, value []byte) error {
	store.Lock()
	store.m[key] = string(value)
	store.Unlock()

	return nil
}

func (store *Store) Del(key string) error {
	delete(store.m, key)
	return nil
}

func (store *Store) Keys() []string {
	store.RLock()
	defer store.RUnlock()
	keys := make([]string, 0)
	for k := range store.m {
		keys = append(keys, k)
	}
	return keys
}

func (store *Store) Close() error {
	return nil
}

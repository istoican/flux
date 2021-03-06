package memory

import "sync"

// Store :
type Store struct {
	sync.RWMutex
	m map[string]string
}

// NewStore :
func NewStore() *Store {
	return &Store{m: make(map[string]string)}
}

// Get :
func (store *Store) Get(key string) ([]byte, error) {
	store.RLock()
	value, _ := store.m[key]
	store.RUnlock()

	return []byte(value), nil
}

// Put :
func (store *Store) Put(key string, value []byte) error {
	store.Lock()
	store.m[key] = string(value)
	store.Unlock()

	return nil
}

// Del :
func (store *Store) Del(key string) error {
	delete(store.m, key)
	return nil
}

// Keys :
func (store *Store) Keys() []string {
	store.RLock()
	defer store.RUnlock()
	keys := make([]string, 0)
	for k := range store.m {
		keys = append(keys, k)
	}
	return keys
}

// Close :
func (store *Store) Close() error {
	return nil
}

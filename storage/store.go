package storage

// An Store is a type that it is used as a backend for storing key value pairs.

type Store interface {
	Get(key string) ([]byte, error)
	Put(key string, value []byte) error
	Del(key string) error
	Keys() []string
	Close() error
}

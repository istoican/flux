package storage

// Store :
type Store interface {
	Get(key string) ([]byte, error)
	Put(key string, value []byte) error
	Del(key string) error
	Keys() []string
	Close() error
}

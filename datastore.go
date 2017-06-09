package flux

// Datastore :
type Datastore interface {
	Get(key string) ([]byte, error)
	Put(key string, value []byte) error
	Del(key string) error
	Close() error
}

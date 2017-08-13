package transport

type Peer interface {
	Get(key string) ([]byte, error)
	Put(key string, value []byte) error
}

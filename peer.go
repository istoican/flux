package flux

// Picker :
type Picker interface {
	Pick(key string) Peer
}

// Peer :
type Peer interface {
	Get(key string) ([]byte, error)
	Put(key string, value []byte) error
}

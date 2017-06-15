package http

// Peer :
type Peer struct {
	address string
}

// Get :
func (peer Peer) Get(key string) ([]byte, error) {
	return []byte(""), nil
}

// Put :
func (peer Peer) Put(key string, value []byte) error {
	return nil
}

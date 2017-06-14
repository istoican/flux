package flux

import (
	"github.com/istoican/flux/consistent"
	"github.com/istoican/flux/storage"
)

// Start :
func Start(config Config) (Flux, error) {
	flux := Flux{
		config:   config,
		listener: newListener(),
		peers:    consistent.New,
	}

	return flux
}

// Flux :
type Node struct {
	name     string
	listener listener
	store    storage.Store
	peers    *consistent.Ring
}

func (node *Node) Lookup(key string) []byte {
	if v, err := node.store.Get(key); err == nil {
		return v
	}
	//peer := node.peers.Get(key)
	//peer.Address

	return ""
}

func (node *Node) Watch(key string) {

}

func (node *Node) Shutdown() error {
	return node.store.Close()
}

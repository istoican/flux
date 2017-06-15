package flux

import (
	"log"

	"github.com/hashicorp/memberlist"
	"github.com/istoican/flux/consistent"
)

// Node :
type Node struct {
	config Config
	peers  *consistent.Ring
	event  listener
}

// Get :
func (node *Node) Get(key string) ([]byte, error) {
	id := node.peers.Get(key).Address
	log.Println("Id: ", id, node.config.ID)
	if id == node.config.ID {
		return node.config.Store.Get(key)
	}
	peer := node.config.Picker.Pick(id)
	return peer.Get(key)
}

// Put :
func (node *Node) Put(key string, value []byte) error {
	id := node.peers.Get(key).Address
	if id == node.config.ID {
		node.event.trigger(key, &Event{Type: "put", Value: value})
		return node.config.Store.Put(key, value)
	}
	peer := node.config.Picker.Pick(id)
	return peer.Put(key, value)
}

// Watch :
func (node *Node) Watch(key string) *Watcher {
	return node.event.watch(key)
}

// Shutdown :
func (node *Node) Shutdown() error {
	return node.config.Store.Close()
}

// NotifyJoin :
func (node *Node) NotifyJoin(n *memberlist.Node) {
	log.Println("Join: ", n.Name)
	node.peers.Add(n.Name)
}

// NotifyLeave :
func (node *Node) NotifyLeave(n *memberlist.Node) {
	node.peers.Remove(n.Address())
}

// NotifyUpdate :
func (node *Node) NotifyUpdate(n *memberlist.Node) {

}

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
	Stats  Stats
}

// Get :
func (node *Node) Get(key string) ([]byte, error) {
	id := node.peers.Get(key).Address
	log.Println("Id: ", id, node.config.ID)
	if id == node.config.ID {
		node.Stats.Reads.Increment()
		return node.config.Store.Get(key)
	}
	peer := node.config.Picker.Pick(id)
	return peer.Get(key)
}

// Put :
func (node *Node) Put(key string, value []byte) error {
	id := node.peers.Get(key).Address
	log.Println("PUT internal id: ", id)
	if id == node.config.ID {
		node.event.trigger(key, &Event{Type: "put", Value: value})
		node.Stats.Inserts.Increment()
		node.Stats.Keys.Increment()

		return node.config.Store.Put(key, value)
	}
	peer := node.config.Picker.Pick(id)
	log.Println("PUT forward to peer: ", peer)
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
	node.peers.Add(n.Name)
	node.config.OnJoin(n.Name)
}

// NotifyLeave :
func (node *Node) NotifyLeave(n *memberlist.Node) {
	node.peers.Remove(n.Address())
	node.config.OnLeave(n.Name)
}

// NotifyUpdate :
func (node *Node) NotifyUpdate(n *memberlist.Node) {

}

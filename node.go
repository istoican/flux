package flux

import (
	"log"
	"time"

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
	//log.Println("Id: ", id, node.config.ID)
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
	//log.Println("PUT internal id: ", id)
	if id == node.config.ID {
		if err := node.config.Store.Put(key, value); err != nil {
			return err
		}
		node.Stats.Inserts.Increment()
		node.Stats.Keys.Set(int64(len(node.config.Store.Keys())))
		node.event.trigger(key, &Event{Type: "put", Value: value})
		return nil
	}
	peer := node.config.Picker.Pick(id)
	//log.Println("PUT forward to peer: ", peer)
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

	go func() {
		peer := node.config.Picker.Pick(n.Name)
		for _, key := range node.config.Store.Keys() {
			id := node.peers.Get(key).Address
			if id != n.Name {
				continue
			}
			if err := node.move(key, peer); err != nil {
				log.Println(err)
			}
			time.Sleep(time.Millisecond)
		}
	}()
}

// NotifyLeave :
func (node *Node) NotifyLeave(n *memberlist.Node) {
	node.peers.Remove(n.Name)
	//log.Println("PEERS2: ", node.peers)
	node.config.OnLeave(n.Name)
}

// NotifyUpdate :
func (node *Node) NotifyUpdate(n *memberlist.Node) {

}

func (node *Node) move(key string, to Peer) error {
	val, err := node.config.Store.Get(key)
	if err != nil {
		return err
	}
	if err := to.Put(key, val); err != nil {
		return err
	}
	if err := node.config.Store.Del(key); err != nil {
		return err
	}
	node.Stats.Keys.Set(int64(len(node.config.Store.Keys())))
	log.Println("MOVED: ", node.Stats.Keys.Value())
	return nil
}

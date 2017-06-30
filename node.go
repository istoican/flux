package flux

import (
	"log"

	"github.com/hashicorp/memberlist"
	"github.com/istoican/flux/consistent"
	"github.com/istoican/flux/storage"
	"github.com/istoican/flux/transport"
)

// Node :
type Node struct {
	addr       string
	store      storage.Store
	ring       *consistent.Ring
	memberlist *memberlist.Memberlist
	metrics    Metrics
	peerFn     func(addr string) transport.Peer
	peers      map[string]transport.Peer
	watchers   map[string][]*Watcher
}

// Addr :
func (n *Node) Local(key string) bool {
	return n.ring.Get(key).Address == n.addr
}

// Addr :
func (n *Node) Peer(key string) (transport.Peer, string) {
	addr := n.ring.Get(key).Address
	peer := n.peers[addr]
	return peer, addr
}

// Get :
func (n *Node) Get(key string) ([]byte, error) {
	addr := n.ring.Get(key).Address

	if addr == n.addr {
		n.metrics.Reads.Increment()
		return n.store.Get(key)
	}
	peer := n.peers[addr]
	return peer.Get(key)
}

// Put :
func (n *Node) Put(key string, value []byte) error {
	addr := n.ring.Get(key).Address

	if addr == n.addr {
		if err := n.store.Put(key, value); err != nil {
			return err
		}
		n.metrics.Inserts.Increment()
		n.metrics.Keys.Set(int64(len(n.store.Keys())))
		n.trigger(key, &Event{Type: "put", Value: string(value)})
		return nil
	}
	peer := n.peers[addr]

	return peer.Put(key, value)
}

// Watch :
func (n *Node) Watch(key string) *Watcher {
	return n.watch(key)
}

// Shutdown :
func (n *Node) Shutdown() error {
	return n.store.Close()
}

// NotifyJoin :
func (n *Node) NotifyJoin(node *memberlist.Node) {
	peer := n.peerFn(node.Name)
	n.ring.Add(node.Name)
	n.peers[node.Name] = peer

	for k := range n.watchers {
		peer2, addr := n.Peer(k)
		if peer2 == peer {
			n.trigger(k, &Event{Type: "moved", Value: addr})
		}
	}
}

// NotifyLeave :
func (n *Node) NotifyLeave(node *memberlist.Node) {
	n.ring.Remove(node.Name)
	delete(n.peers, node.Name)
}

// NotifyUpdate :
func (n *Node) NotifyUpdate(node *memberlist.Node) {

}

func (n *Node) rebalance() {
	for _, key := range n.store.Keys() {
		addr := n.ring.Get(key).Address
		if addr == n.addr {
			continue
		}
		peer := n.peers[addr]
		if err := n.move(key, peer); err != nil {
			log.Println(err)
		}
	}
}

func (n *Node) move(key string, to transport.Peer) error {
	val, err := n.store.Get(key)
	if err != nil {
		return err
	}
	if err := to.Put(key, val); err != nil {
		return err
	}
	if err := n.store.Del(key); err != nil {
		return err
	}
	n.metrics.Keys.Set(int64(len(n.store.Keys())))
	return nil
}

// Metrics :
func (n *Node) Metrics() interface{} {
	members := make(map[string]string)
	for _, v := range n.memberlist.Members() {
		members[v.Name] = v.Addr.String()
	}
	return map[string]interface{}{
		"stats":   n.metrics,
		"members": members,
	}
}

// Join :
func (n *Node) Join(address string) error {
	if address == "" {
		return nil
	}
	_, err := n.memberlist.Join([]string{address})

	return err
}

func (n *Node) watch(path string) *Watcher {
	w := &Watcher{Channel: make(chan *Event, 100)}

	if _, ok := n.watchers[path]; !ok {
		n.watchers[path] = make([]*Watcher, 0)
	}

	n.watchers[path] = append(n.watchers[path], w)

	w.Remove = func() {
	}

	return w
}

func (n *Node) trigger(path string, e *Event) {
	watchers, ok := n.watchers[path]
	if !ok {
		return
	}

	for _, w := range watchers {
		w.notify(e)
	}
}

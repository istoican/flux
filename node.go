package flux

import (
	"log"

	"sync"

	"time"

	"github.com/hashicorp/memberlist"
	"github.com/istoican/flux/consistent"
	"github.com/istoican/flux/storage"
	"github.com/istoican/flux/transport"
)

// Represents a node in the distributed system. It manges key distribution using a consistent hash ring
// Interacts with other nodes using the memberlist package.
type Node struct {
	mu         sync.RWMutex
	addr       string
	store      storage.Store
	ring       *consistent.Ring
	memberlist *memberlist.Memberlist
	metrics    Metrics
	peerFn     func(addr string) transport.Peer
	peers      map[string]transport.Peer
	watchers   map[string][]*Watcher
}

// Checks if a given key belongs to the local node based on the consistent hash ring.
func (n *Node) Local(key string) bool {
	return n.ring.Get(key).Address == n.addr
}

// Retrieves the peer responsible for a given key and returns both the peer and its address.
func (n *Node) Peer(key string) (transport.Peer, string) {
	n.mu.RLock()
	defer n.mu.RUnlock()

	addr := n.ring.Get(key).Address
	peer := n.peers[addr]
	return peer, addr
}

// Retrieves the value for a given key. If the key is not local, it forwards the request to the appropriate peer.
func (n *Node) Get(key string) ([]byte, error) {
	n.mu.Lock()
	defer n.mu.Unlock()

	addr := n.ring.Get(key).Address

	if addr == n.addr {
		n.metrics.Reads.Increment()
		return n.store.Get(key)
	}
	peer := n.peers[addr]
	return peer.Get(key)
}

// Stores a key-value pair. It the key is not local, it forwards the request to the appropriate peer.
func (n *Node) Put(key string, value []byte) error {
	n.mu.Lock()
	defer n.mu.Unlock()

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

// Sets up a watcher for a given key.
func (n *Node) Watch(key string) *Watcher {
	return n.watch(key)
}

// Close the node's attached storage.
func (n *Node) Shutdown() error {
	return n.store.Close()
}

// Handle node membership changes in the cluster.
func (n *Node) NotifyJoin(node *memberlist.Node) {
	n.mu.Lock()
	defer n.mu.Unlock()

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

func (n *Node) NotifyLeave(node *memberlist.Node) {
	n.mu.Lock()
	defer n.mu.Unlock()

	n.ring.Remove(node.Name)
	delete(n.peers, node.Name)
}

func (n *Node) NotifyUpdate(node *memberlist.Node) {

}

// Iterate over local keys and moves unassigned keys to the correct node.
func (n *Node) rebalance() {
	for _, key := range n.store.Keys() {
		n.mu.Lock()
		addr := n.ring.Get(key).Address
		if addr == n.addr {
			n.mu.Unlock()
			continue
		}
		peer := n.peers[addr]
		if err := n.move(key, peer); err != nil {
			log.Println(err)
		}
		n.mu.Unlock()
		time.Sleep(1 * time.Millisecond)
	}
}

// Move a key from the local node to another node.
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

// Return statistics and membership information.
func (n *Node) Metrics() interface{} {
	n.mu.RLock()
	defer n.mu.RUnlock()

	members := make(map[string]string)
	for _, v := range n.memberlist.Members() {
		members[v.Name] = v.Addr.String()
	}
	return map[string]interface{}{
		"stats":   n.metrics,
		"members": members,
	}
}

// Add an address to the cluster.
func (n *Node) Join(address string) error {
	if address == "" {
		return nil
	}
	_, err := n.memberlist.Join([]string{address})

	return err
}

// Set up a watcher for a given path.
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

// Notifies watchers for a given path.
func (n *Node) trigger(path string, e *Event) {
	watchers, ok := n.watchers[path]
	if !ok {
		return
	}

	for _, w := range watchers {
		w.notify(e)
	}
}

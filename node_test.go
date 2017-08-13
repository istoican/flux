package flux

import (
	"errors"
	"testing"

	"github.com/hashicorp/memberlist"
	"github.com/istoican/flux/consistent"
	"github.com/istoican/flux/consistent/hash"
	"github.com/istoican/flux/storage/memory"
	"github.com/istoican/flux/transport"
)

type MockPeer struct {
	store map[string][]byte
}

func (p *MockPeer) Get(key string) ([]byte, error) {
	if val, ok := p.store[key]; ok {
		return val, nil
	}
	return nil, errors.New("key not found")
}

func (p *MockPeer) Put(key string, value []byte) error {
	p.store[key] = value
	return nil
}

func (p *MockPeer) Del(key string) error {
	delete(p.store, key)
	return nil
}

func TestNodeLocal(t *testing.T) {
	store := memory.NewStore()
	ring := consistent.New(hash.MD5)

	n := &Node{
		addr:  "node1",
		store: store,
		ring:  ring,
		peers: make(map[string]transport.Peer),
	}

	n.ring.Add("node1")

	if !n.Local("key1") {
		t.Errorf("expected key to be local to node1")
	}

	n.ring.Add("node2")

	if n.Local("key1") && n.ring.Get("key1").Address != "node1" {
		t.Errorf("expected key to not be local to node1")
	}
}

func TestNodePeer(t *testing.T) {
	store := memory.NewStore()
	ring := consistent.New(hash.MD5)

	n := &Node{
		addr:  "node1",
		store: store,
		ring:  ring,
		peers: make(map[string]transport.Peer),
		peerFn: func(addr string) transport.Peer {
			return &MockPeer{store: make(map[string][]byte)}
		},
	}

	n.ring.Add("node1")
	n.ring.Add("node2")
	n.peers["node1"] = &MockPeer{store: make(map[string][]byte)}
	n.peers["node2"] = &MockPeer{store: make(map[string][]byte)}

	peer, addr := n.Peer("key1")
	if peer == nil || addr != "node1" && addr != "node2" {
		t.Errorf("expected to find a peer for key1")
	}
}

func TestNodeGetPut(t *testing.T) {
	store := memory.NewStore()
	ring := consistent.New(hash.MD5)
	n := &Node{
		addr:  "node1",
		store: store,
		ring:  ring,
		peers: make(map[string]transport.Peer),
		peerFn: func(addr string) transport.Peer {
			return &MockPeer{store: make(map[string][]byte)}
		},
	}

	n.ring.Add("node1")
	n.peers["node1"] = &MockPeer{store: make(map[string][]byte)}

	err := n.Put("key1", []byte("value1"))
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	val, err := n.Get("key1")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if string(val) != "value1" {
		t.Errorf("expected 'value1', got %s", val)
	}
}

func TestNodeNotifyJoinAndLeave(t *testing.T) {
	store := memory.NewStore()
	ring := consistent.New(hash.MD5)
	n := &Node{
		addr:  "node1",
		store: store,
		ring:  ring,
		peers: make(map[string]transport.Peer),
		peerFn: func(addr string) transport.Peer {
			return &MockPeer{store: make(map[string][]byte)}
		},
	}

	node := &memberlist.Node{Name: "node2"}
	n.NotifyJoin(node)

	if len(n.peers) != 1 {
		t.Errorf("expected 1 peer, got %d", len(n.peers))
	}

	n.NotifyLeave(node)
	if len(n.peers) != 0 {
		t.Errorf("expected 0 peers, got %d", len(n.peers))
	}
}

func TestNodeRebalance(t *testing.T) {
	store := memory.NewStore()
	ring := consistent.New(hash.MD5)
	n := &Node{
		addr:  "node1",
		store: store,
		ring:  ring,
		peers: make(map[string]transport.Peer),
		peerFn: func(addr string) transport.Peer {
			return &MockPeer{store: make(map[string][]byte)}
		},
	}

	n.ring.Add("node1")
	n.peers["node1"] = &MockPeer{store: make(map[string][]byte)}

	n.Put("key1", []byte("value1"))

	node2Peer := &MockPeer{store: make(map[string][]byte)}
	n.peers["node2"] = node2Peer
	n.ring.Add("node2")

	n.rebalance()

	if _, err := node2Peer.Get("key1"); err != nil {
		t.Errorf("expected key to be moved to node2")
	}
}

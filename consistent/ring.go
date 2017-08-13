package consistent

import (
	"sort"
	"strconv"
)

const (
	replicas = 256
)

// HashFn is responsible for generating unsigned, 32-bit hash of provided string.
// It should minimize collisions (generating same hash for different string
type HashFn func(string) uint32

// Ring provides an implementation of a ring hash.
type Ring struct {
	nodes  Nodes
	hashFn HashFn
}

// Add a new node to the ring
func (r *Ring) Add(address string) {
	for i := 0; i < replicas; i++ {
		node := Node{
			Address: address,
			hash:    r.hashFn(strconv.Itoa(i) + address),
		}
		r.nodes = append(r.nodes, node)
	}
	sort.Sort(r.nodes)
}

// Remove a node from the ring
func (r *Ring) Remove(address string) {
	nodes := r.nodes[:0]
	for _, n := range r.nodes {
		if n.Address != address {
			nodes = append(nodes, n)
		}
	}
	r.nodes = nodes
}

// Giving a key it returns an node
func (r *Ring) Get(key string) Node {
	hash := r.hashFn(key)
	f := func(i int) bool {
		return r.nodes[i].hash >= hash
	}
	i := sort.Search(r.nodes.Len(), f)
	if i >= r.nodes.Len() {
		i = 0
	}

	return r.nodes[i]
}

package consistent

import (
	"sort"
)

const (
	replicas = 256
)

// HashFn :
type HashFn func(string) uint32

// Ring :
type Ring struct {
	nodes  Nodes
	hashFn HashFn
}

// Add : Add a new node to the ring
func (r *Ring) Add(address string) {
	for i := 0; i < replicas; i++ {
		node := Node{
			Address: address,
			hash:    r.hashFn(string(i) + address),
		}
		r.nodes = append(r.nodes, node)
	}
	sort.Sort(r.nodes)
}

// Remove :
func (r *Ring) Remove(address string) {
	nodes := r.nodes[:0]
	for _, n := range r.nodes {
		if n.Address != address {
			nodes = append(nodes, n)
		}
	}
	r.nodes = nodes
}

// Get :
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

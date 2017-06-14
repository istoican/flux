package consistent

import (
	"hash/crc32"
	"sort"
)

// Ring :
type Ring struct {
	nodes Nodes
}

// Add : Add anew node to the ring
func (r *Ring) Add(address string) {
	node := Node{
		Address: address,
		hash:    crc32.ChecksumIEEE([]byte(address)),
	}
	nodes := append(r.nodes, node)
	sort.Sort(nodes)
	r.nodes = nodes
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
	f := func(i int) bool {
		return r.nodes[i].hash >= crc32.ChecksumIEEE([]byte(key))
	}
	i := sort.Search(r.nodes.Len(), f)
	if i >= r.nodes.Len() {
		i = 0
	}
	return r.nodes[i]
}

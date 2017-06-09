package consistent

import (
	"hash/crc32"
	"sort"
)

// Ring :
type Ring struct {
	nodes Nodes
}

// Add :
func (r *Ring) Add(id string) {
	node := Node{
		ID:   id,
		hash: crc32.ChecksumIEEE([]byte(id)),
	}
	r.add(node)
}

func (r *Ring) add(node Node) {
	nodes := append(r.nodes, node)
	sort.Sort(nodes)
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

package consistent

import (
	"sort"
)

const (
	replicas = 256
)

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

// HashFn :
type HashFn func(string) uint32

/*
func hash(key string) uint32 {
	if len(key) < 64 {
		var scratch [64]byte
		copy(scratch[:], key)
		return crc32.ChecksumIEEE(scratch[:len(key)])
	}
	return crc32.ChecksumIEEE([]byte(key))
}
*/

//func MD5Hhash(key string) string {
//	m := md5.New()
//	m.Write([]byte(key))
//	return string(m.Sum(nil))
//}

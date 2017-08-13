package consistent

// Package consistent provides a consistent hashing function.

func New(fn HashFn) *Ring {
	nodes := []Node{}
	return &Ring{
		nodes:  nodes,
		hashFn: fn,
	}
}

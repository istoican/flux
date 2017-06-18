package consistent

// New :
func New(fn HashFn) *Ring {
	nodes := []Node{}
	return &Ring{
		nodes:  nodes,
		hashFn: fn,
	}
}

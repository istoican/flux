package consistent

// New :
func New() *Ring {
	nodes := []Node{}
	return &Ring{
		nodes: nodes,
	}
}

package consistent

// Node :
type Node struct {
	Address string
	hash    uint32
}

// Nodes :
type Nodes []Node

// Len :
func (n Nodes) Len() int {
	return len(n)
}

// Less :
func (n Nodes) Less(i, j int) bool {
	return n[i].hash < n[j].hash
}

// Swap :
func (n Nodes) Swap(i, j int) {
	n[i], n[j] = n[j], n[i]
}

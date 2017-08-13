package consistent

type Node struct {
	Address string
	hash    uint32
}

type Nodes []Node

func (n Nodes) Len() int {
	return len(n)
}

func (n Nodes) Less(i, j int) bool {
	return n[i].hash < n[j].hash
}

func (n Nodes) Swap(i, j int) {
	n[i], n[j] = n[j], n[i]
}

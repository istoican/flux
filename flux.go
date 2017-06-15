package flux

import (
	"github.com/hashicorp/memberlist"
	"github.com/istoican/flux/consistent"
)

var (
	node       Node
	membership *memberlist.Memberlist
)

// Start :
func Start(config Config) error {
	n := Node{
		config: config,
		event:  newListener(),
		peers:  consistent.New(),
		Stats:  Stats{},
	}

	memberlistConfig := memberlist.DefaultLocalConfig()
	memberlistConfig.Events = &n

	l, err := memberlist.Create(memberlistConfig)
	if err != nil {
		return err
	}
	node = n
	membership = l
	return nil
}

// Get :
func Get(key string) ([]byte, error) {
	return node.Get(key)
}

// Put :
func Put(key string, value []byte) error {
	return node.Put(key, value)
}

// Watch :
func Watch(key string) *Watcher {
	return node.Watch(key)
}

// Join :
func Join(address string) error {
	if address == "" {
		return nil
	}
	_, err := membership.Join([]string{address})

	return err
}

package flux

import (
	"github.com/hashicorp/memberlist"
	"github.com/istoican/flux/consistent"
)

var (
	node Node
)

// Start :
func Start(config Config) error {
	n := Node{
		config: config,
		event:  newListener(),
		peers:  consistent.New(),
	}

	memberlistConfig := memberlist.DefaultLocalConfig()
	memberlistConfig.Events = &n

	_, err := memberlist.Create(memberlistConfig)
	if err != nil {
		return err
	}
	node = n
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

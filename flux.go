package flux

import (
	"time"

	"github.com/hashicorp/memberlist"
	"github.com/istoican/flux/consistent"
	"github.com/istoican/flux/transport"
)

// Initializes a new Node with the provided configuration.
// It sets up the consistent hash ring, initializes peer connections,
// configures the memberlist for managing cluster membership,
// and starts a goroutine for periodically rebalancing local stored keys.
func New(config Config) (*Node, error) {
	n := &Node{
		addr:     config.Addr,
		store:    config.Store,
		ring:     consistent.New(config.HashFn),
		peers:    make(map[string]transport.Peer),
		peerFn:   config.PeerFn,
		metrics:  Metrics{},
		watchers: make(map[string][]*Watcher),
	}

	memberlistConfig := memberlist.DefaultLocalConfig()
	memberlistConfig.Events = n

	memberlist, err := memberlist.Create(memberlistConfig)
	if err != nil {
		return nil, err
	}
	n.memberlist = memberlist

	// periodically rebalance local stored keys
	go func() {
		for {
			n.rebalance()
			time.Sleep(1 * time.Second)
		}
	}()

	return n, nil
}

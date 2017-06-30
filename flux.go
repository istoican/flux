package flux

import (
	"time"

	"github.com/hashicorp/memberlist"
	"github.com/istoican/flux/consistent"
	"github.com/istoican/flux/transport"
)

// New :
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

	go func() {
		for {
			n.rebalance()
			time.Sleep(1 * time.Second)
		}
	}()
	return n, nil
}

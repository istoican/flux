package flux

import (
	"os"

	"github.com/istoican/flux/consistent"
	"github.com/istoican/flux/consistent/hash"
	"github.com/istoican/flux/storage"
	"github.com/istoican/flux/storage/memory"
	"github.com/istoican/flux/transport"
	"github.com/istoican/flux/transport/http/peer"
)

// Config :
type Config struct {
	Addr    string
	Store   storage.Store
	OnJoin  func(id string)
	OnLeave func(id string)
	PeerFn  func(string) transport.Peer
	HashFn  consistent.HashFn
}

// DefaultConfig :
func DefaultConfig() Config {
	hostname, _ := os.Hostname()

	return Config{
		Addr:   hostname,
		Store:  memory.NewStore(),
		HashFn: hash.CRC32,
		PeerFn: peer.New,
	}
}

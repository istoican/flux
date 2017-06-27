package flux

import (
	"os"

	"github.com/istoican/flux/consistent"
	"github.com/istoican/flux/consistent/hash"
	"github.com/istoican/flux/storage"
	"github.com/istoican/flux/storage/memory"
)

// Config :
type Config struct {
	ID      string
	Store   storage.Store
	OnJoin  func(id string)
	OnLeave func(id string)
	Picker  Picker
	HashFn  consistent.HashFn
}

// DefaultConfig :
func DefaultConfig() Config {
	hostname, _ := os.Hostname()

	return Config{
		ID:     hostname,
		Store:  memory.NewStore(),
		HashFn: hash.CRC32,
	}
}

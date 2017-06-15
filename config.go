package flux

import (
	"github.com/istoican/flux/storage"
)

// Config :
type Config struct {
	ID      string
	Store   storage.Store
	OnJoin  func(id string)
	OnLeave func(id string)
	Picker  Picker
}

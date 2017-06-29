package flux

import (
	"sync/atomic"
)

// Stats :
type Stats struct {
	Keys      Int
	Deletions Int
	Inserts   Int
	Reads     Int
}

// Int :
type Int int64

// Increment :
func (c *Int) Increment() {
	atomic.AddInt64((*int64)(c), 1)
}

// Decrement :
func (c *Int) Decrement() {
	atomic.AddInt64((*int64)(c), -1)
}

// Set :
func (c *Int) Set(v int64) {
	atomic.StoreInt64((*int64)(c), v)
}

package flux

import (
	"sync/atomic"
)

// Metrics holds the counters for various operations in the node such as key inserts, deletions, reads, and total number of keys.
type Metrics struct {
	Keys      Int
	Deletions Int
	Inserts   Int
	Reads     Int
}

type Int int64

func (c *Int) Increment() {
	atomic.AddInt64((*int64)(c), 1)
}

func (c *Int) Decrement() {
	atomic.AddInt64((*int64)(c), -1)
}

func (c *Int) Set(v int64) {
	atomic.StoreInt64((*int64)(c), v)
}

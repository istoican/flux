package flux

import (
	"strconv"
	"sync/atomic"
)

// Statistics :
type Stats struct {
	Keys      Counter
	Deletions Counter
	Inserts   Counter
	Reads     Counter
}

type Counter int64

func (c *Counter) Increment() {
	atomic.AddInt64((*int64)(c), 1)
}

func (c *Counter) Decrement() {
	atomic.AddInt64((*int64)(c), -1)
}

func (c *Counter) Get() int64 {
	return atomic.LoadInt64((*int64)(c))
}

func (c *Counter) String() string {
	return strconv.FormatInt(c.Get(), 10)
}

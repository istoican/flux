package flux

import (
	"testing"
	"time"
)

func TestWatcher(t *testing.T) {
	ch := make(chan *Event)

	watcher := Watcher{
		Channel: ch,
	}

	go func() {
		watcher.notify(&Event{})
	}()

	select {
	case <-ch:
	case <-time.After(5 * time.Second):
		t.Fatalf("watcher timed out")
	}
}

package flux

import (
	"container/list"
	"log"
)

type listener struct {
	watchers map[string]*list.List
}

// Watcher :
type Watcher struct {
	Channel chan *Event
	Remove  func()
}

// Event :
type Event struct {
	Action string
	Value  interface{}
}

func newListener() listener {
	return listener{watchers: make(map[string]*list.List)}
}

func (l *listener) watch(path string) *Watcher {
	w := &Watcher{Channel: make(chan *Event, 100)}

	if _, ok := l.watchers[path]; !ok {
		l.watchers[path] = list.New()
	}

	e := l.watchers[path].PushBack(w)

	w.Remove = func() {
		l.watchers[path].Remove(e)
	}

	return w
}

func (l *listener) trigger(path string, e *Event) {
	elements, ok := l.watchers[path]
	if !ok {
		return
	}

	curr := elements.Front()
	for curr != nil {
		w, err := curr.Value.(*Watcher)

		if err {
			log.Println(err)
		}
		w.notify(e)
		curr = curr.Next()
	}
}

func (w *Watcher) notify(e *Event) {
	w.Channel <- e
}

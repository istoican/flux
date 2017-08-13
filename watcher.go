package flux

type Event struct {
	Type  string
	Value interface{}
}

// Watcher implements a simple event notification bus system
type Watcher struct {
	Channel chan *Event
	Remove  func()
}

func (w *Watcher) notify(e *Event) {
	w.Channel <- e
}

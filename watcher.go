package flux

// Watcher :
type Watcher struct {
	Channel chan *Event
	Remove  func()
}

func (w *Watcher) notify(e *Event) {
	w.Channel <- e
}

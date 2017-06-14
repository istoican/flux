package flux

import (
	"encoding/json"
)

// Get :
func (server *Server) Get(key string) (interface{}, error) {
	data, err := server.store.Get(key)
	if err != nil {
		return nil, err
	}

	var val map[string]interface{}

	if err := json.Unmarshal(data, &val); err != nil {
		return nil, err
	}

	return val, nil
}

// Close :

// Delete :
func (server *Server) Delete(key string) error {
	server.event.trigger(key, &Event{Action: "delete"})
	return server.store.Del(key)
}

// Watch :
func (server *Server) Watch(key string) *Watcher {
	return server.event.watch(key)
}

// Put :
func (server *Server) Put(key string, value interface{}) error {
	val, err := json.Marshal(value)
	if err != nil {
		return err
	}
	server.event.trigger(key, &Event{Action: "change", Value: value})

	return server.store.Put(key, val)
}

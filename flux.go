package flux

import (
	"encoding/json"

	"github.com/istoican/flux/backend"
)

type DB struct {
	event listener
	store Datastore
}

func (db *DB) Get(key string) (interface{}, error) {
	data, err := db.store.Get(key)
	if err != nil {
		return nil, err
	}

	var val map[string]interface{}

	if err := json.Unmarshal(data, &val); err != nil {
		return nil, err
	}

	return val, nil
}

func (db *DB) Close() error {
	return db.store.Close()
}

func (db *DB) Delete(key string) error {
	db.event.trigger(key, &Event{Action: "delete"})
	return db.store.Del(key)
}

func (db *DB) Watch(key string) *Watcher {
	return db.event.watch(key)
}

func (db *DB) Put(key string, value interface{}) error {
	val, err := json.Marshal(value)
	if err != nil {
		return err
	}
	db.event.trigger(key, &Event{Action: "change", Value: value})

	return db.store.Put(key, val)
}

func Open(file string) (*DB, error) {
	b, err := backend.NewGoLevelDB(file)
	if err != nil {
		return nil, err
	}

	db := &DB{
		event: newListener(),
		store: b,
	}
	return db, nil
}

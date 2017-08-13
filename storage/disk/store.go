package disk

import (
	"github.com/syndtr/goleveldb/leveldb"
)

// Disk provides an persistent implementation of Store interface. It is backed by levelDB.

type GoLevelDB struct {
	backend *leveldb.DB
}

func NewStore(path string) (*GoLevelDB, error) {
	db, err := leveldb.OpenFile(path, nil)
	if err != nil {
		return nil, err
	}
	return &GoLevelDB{backend: db}, nil
}

func (db *GoLevelDB) Get(key string) ([]byte, error) {
	return db.backend.Get([]byte(key), nil)
}

func (db *GoLevelDB) Put(key string, value []byte) error {
	return db.backend.Put([]byte(key), value, nil)
}

func (db *GoLevelDB) Del(key string) error {
	return db.backend.Delete([]byte(key), nil)
}

func (db *GoLevelDB) Keys() []string {
	it := db.backend.NewIterator(nil, nil)
	keys := make([]string, 0)
	for it.Next() {
		keys = append(keys, string(it.Key()))
	}
	return keys
}

func (db *GoLevelDB) Close() error {
	return db.backend.Close()
}

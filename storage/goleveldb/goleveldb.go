package backend

import (
	"github.com/syndtr/goleveldb/leveldb"
)

type GoLevelDB struct {
	backend *leveldb.DB
}

func NewGoLevelDB(path string) (*GoLevelDB, error) {
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

func (db *GoLevelDB) Close() error {
	return db.backend.Close()
}

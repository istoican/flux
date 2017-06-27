package disk

import (
	"github.com/syndtr/goleveldb/leveldb"
)

// GoLevelDB :
type GoLevelDB struct {
	backend *leveldb.DB
}

// NewStore :
func NewStore(path string) (*GoLevelDB, error) {
	db, err := leveldb.OpenFile(path, nil)
	if err != nil {
		return nil, err
	}
	return &GoLevelDB{backend: db}, nil
}

// Get :
func (db *GoLevelDB) Get(key string) ([]byte, error) {
	return db.backend.Get([]byte(key), nil)
}

// Put :
func (db *GoLevelDB) Put(key string, value []byte) error {
	return db.backend.Put([]byte(key), value, nil)
}

// Del :
func (db *GoLevelDB) Del(key string) error {
	return db.backend.Delete([]byte(key), nil)
}

// Close :
func (db *GoLevelDB) Close() error {
	return db.backend.Close()
}

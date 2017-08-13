package hash

import (
	"crypto/md5"
	"math/big"
)

func MD5(key string) (out uint32) {
	m := md5.New()
	m.Write([]byte(key))
	hash := m.Sum(nil)

	return uint32(new(big.Int).SetBytes(hash).Uint64())
}

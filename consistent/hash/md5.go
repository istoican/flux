package hash

import (
	"crypto/md5"
	"math/big"
)

// MD5 :
func MD5(key string) (out uint32) {
	m := md5.New()
	m.Write([]byte(key))
	hash := m.Sum(nil)

	return uint32(new(big.Int).SetBytes(hash).Uint64())
	//for i, b := range hash {
	//	shift := uint32((16 - i - 1) * 8)
	//
	//		out |= uint32(b) << shift
	//	}
	//return
}

package hash

import (
	"crypto/md5"
)

// MD5 :
func MD5(key string) (out uint32) {
	m := md5.New()
	m.Write([]byte(key))
	hash := string(m.Sum(nil))

	for i, b := range hash {
		shift := uint32((16 - i - 1) * 8)

		out |= uint32(b) << shift
	}
	return
}

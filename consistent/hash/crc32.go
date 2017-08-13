package hash

import (
	"hash/crc32"
)

func CRC32(key string) uint32 {
	return crc32.ChecksumIEEE([]byte(key))
}

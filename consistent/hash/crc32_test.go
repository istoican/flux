package hash

import (
	"testing"
)

func TestCRC32(t *testing.T) {
	testCases := map[string]uint32{
		"2": 450215437,
		"4": 4088798008,
		"6": 498629140,
	}

	for k, v := range testCases {
		if CRC32(k) != v {
			t.Errorf("Asking for %s, should have yielded %v", k, v)
		}
	}
}

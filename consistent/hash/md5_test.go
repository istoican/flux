package hash

import (
	"testing"
)

func TestMD5(t *testing.T) {
	testCases := map[string]uint32{
		"2": 3423897132,
		"4": 1967264300,
		"6": 2125574876,
	}

	for k, v := range testCases {
		if MD5(k) != v {
			t.Errorf("Asking for %s, should have yielded %v", k, v)
		}
	}
}

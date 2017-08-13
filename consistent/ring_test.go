package consistent

import (
	"testing"

	"github.com/istoican/flux/consistent/hash"
)

func TestRing(t *testing.T) {
	ring := Ring{hashFn: hash.MD5}

	nodes := []string{"2", "4", "6"}

	for _, n := range nodes {
		ring.Add(n)
	}

	testCases := map[string]string{
		"9":  "6",
		"10": "4",
		"11": "6",
		"12": "2",
		"13": "6",
		"14": "4",
	}

	for k, v := range testCases {
		n := ring.Get(k)
		if n.Address != v {
			t.Errorf("Asking for %s, should have yielded %v", k, v)
		}
	}
}

package memory

import (
	"bytes"
	"testing"
)

func TestStore(t *testing.T) {
	store := NewStore()

	testCases := map[string][]byte{
		"1": []byte("2"),
		"2": []byte("4"),
		"3": []byte("6"),
	}

	for k, v := range testCases {
		if err := store.Put(k, v); err != nil {
			t.Errorf(" Error adding value to storer %s", err)
		}
	}

	for k, v := range testCases {
		v2, err := store.Get(k)
		if err != nil {
			t.Errorf(" Error adding value to storer %s", err)
		}

		if !bytes.Equal(v, v2) {
			t.Errorf("Asking for %s, should have yielded %v", k, v)
		}
	}
}

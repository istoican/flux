package peer

import (
	"bytes"
	"io/ioutil"
	"net/http"

	"github.com/istoican/flux/transport"
)

// Peer :
type Peer struct {
	address string
}

// New :
func New(addr string) transport.Peer {
	return Peer{addr}
}

func (peer Peer) String() string {
	return peer.address
}

// Get :
func (peer Peer) Get(key string) ([]byte, error) {
	resp, err := http.Get("http://" + peer.address + ":8080/" + key)
	if err != nil {
		return []byte(""), err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []byte(""), err
	}
	return body, nil
}

// Put :
func (peer Peer) Put(key string, value []byte) error {
	//log.Println("http PUT: ", peer.address)
	buf := bytes.NewReader(value)
	resp, err := http.Post("http://"+peer.address+":8080/"+key, "text/plain", buf)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}

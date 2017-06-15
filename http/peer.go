package http

import (
	"bytes"
	"io/ioutil"
	"net/http"
)

// Peer :
type Peer struct {
	address string
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
	_, err := http.Post("http://"+peer.address+":8080/"+key, "text/plain", buf)
	if err != nil {
		return err
	}
	return nil
}

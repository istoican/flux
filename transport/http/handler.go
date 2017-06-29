package http

import (
	"encoding/json"
	"expvar"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"sync"

	"github.com/istoican/flux"
)

var (
	handler *httpHandler
)

// Join :
func Join(id string) {
	log.Println("JOIN: ", id)
	handler.mu.Lock()
	defer handler.mu.Unlock()
	handler.peers[id] = &Peer{address: id}
}

// Leave :
func Leave(id string) {
	handler.mu.Lock()
	defer handler.mu.Unlock()
	log.Println("DEL: ", id)
	delete(handler.peers, id)
	log.Println("PEERS: ", handler.peers)
}

// Handler :
func Handler() flux.Picker {
	return handler
}

func init() {
	handler = &httpHandler{
		peers: make(map[string]*Peer),
	}

	http.Handle("/", handler)
	http.HandleFunc("/debug/nodes", nodesHandler)
	expvar.Publish("flux", expvar.Func(flux.Info))
}

type httpHandler struct {
	mu    sync.Mutex
	peers map[string]*Peer
}

func (handler *httpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Path[1:]

	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case "GET":
		if r.FormValue("watch") == "" {
			value, err := flux.Get(key)
			log.Printf("GET(%s)", key)
			response(w, string(value), err)
			return
		}

		watcher := flux.Watch(key)

		defer func() {
			if watcher.Remove != nil {
				watcher.Remove()
			}
		}()

		ch := watcher.Channel

		for {
			select {
			case e, ok := <-ch:
				if !ok {
					return
				}
				response(w, fmt.Sprintf("%s", e.Value), nil)
				w.(http.Flusher).Flush()
			}
		}
	case "POST":
		body, _ := ioutil.ReadAll(r.Body)
		v := string(body)
		log.Printf("PUT(%s, %s)", key, v)
		err := flux.Put(key, []byte(v))
		response(w, v, err)
	}
}

// Pick :
func (handler *httpHandler) Pick(key string) flux.Peer {
	handler.mu.Lock()
	defer handler.mu.Unlock()
	peer, _ := handler.peers[key]
	return peer
}

func response(w io.Writer, value interface{}, err error) {
	if err := json.NewEncoder(w).Encode(value); err != nil {
		log.Printf("error %v \n", err)
	}
}

func nodesHandler(w http.ResponseWriter, r *http.Request) {
	nodes := make(map[string]string)

	for _, n := range flux.Nodes() {
		nodes[n.Name] = n.Addr.String()
	}

	if err := json.NewEncoder(w).Encode(nodes); err != nil {
		log.Printf("error %v \n", err)
	}
}

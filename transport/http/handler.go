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

	"os"

	"time"

	"github.com/istoican/flux"
	"github.com/istoican/flux/storage/memory"
)

var (
	stats = expvar.NewMap("flux")
)

func init() {
	handler := &Handler{
		peers: make(map[string]*Peer),
	}

	onJoin := func(id string) {
		//log.Println("http JOIN: ", id)
		handler.mu.Lock()
		defer handler.mu.Unlock()
		handler.peers[id] = &Peer{address: id}
	}

	onLeave := func(id string) {
		handler.mu.Lock()
		defer handler.mu.Unlock()
		delete(handler.peers, id)
	}

	hostname, _ := os.Hostname()

	config := flux.Config{
		ID:      hostname,
		Store:   memory.NewStore(),
		OnJoin:  onJoin,
		OnLeave: onLeave,
		Picker:  handler,
	}
	flux.Start(config)

	http.Handle("/", handler)
	http.HandleFunc("/debug/nodes", nodesHandler)

	go func() {
		for {
			s := flux.Info()

			stats.Set("keys", &s.Keys)
			stats.Set("inserts", &s.Inserts)
			stats.Set("deletions", &s.Deletions)
			stats.Set("reads", &s.Reads)

			time.Sleep(10 * time.Millisecond)
		}
	}()
}

// Handler :
type Handler struct {
	mu    sync.Mutex
	peers map[string]*Peer
}

func (handler *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Path[1:]

	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case "GET":
		if r.FormValue("watch") == "" {
			value, err := flux.Get(key)
			log.Printf("GET(%s)", value)
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
		//case "DELETE":
		//	err := flux.Delete(key)
		//	response(w, nil, err)
	}
}

// Pick :
func (handler *Handler) Pick(key string) flux.Peer {
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

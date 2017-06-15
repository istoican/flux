package http

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"sync"

	"os"

	"github.com/istoican/flux"
	"github.com/istoican/flux/storage/memory"
)

func init() {
	handler := &Handler{}

	onJoin := func(id string) {
		handler.mu.Lock()
		defer handler.mu.Unlock()
		handler.peers[id] = Peer{id}
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
}

// Handler :
type Handler struct {
	mu    sync.Mutex
	peers map[string]Peer
}

func (handler *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Path[1:]

	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case "GET":
		if r.FormValue("watch") == "" {
			value, err := flux.Get(key)
			log.Println("GET: ", string(value), err)
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
		log.Println("PUT: ", v)
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

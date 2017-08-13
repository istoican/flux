package handler

import (
	"encoding/json"
	"expvar"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/istoican/flux"
)

func New(node *flux.Node) http.Handler {
	expvar.Publish("flux", expvar.Func(node.Metrics))
	return &server{node}
}

type server struct {
	node *flux.Node
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Path[1:]

	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case "GET":
		if r.FormValue("watch") == "" {
			value, err := s.node.Get(key)
			log.Printf("GET(%s)", key)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Write(value)
			return
		}

		if !s.node.Local(key) {
			_, addr := s.node.Peer(key)
			if err := json.NewEncoder(w).Encode(flux.Event{Type: "moved", Value: addr}); err != nil {
				log.Println(err)
			}
			w.(http.Flusher).Flush()
		}
		watcher := s.node.Watch(key)

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
				if err := json.NewEncoder(w).Encode(e); err != nil {
					log.Printf("error %v \n", err)
				}
				w.(http.Flusher).Flush()
			}
		}
	case "POST":
		body, _ := ioutil.ReadAll(r.Body)
		log.Printf("PUT(%s, %s)", key, body)
		err := s.node.Put(key, body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Write(body)
	}
}

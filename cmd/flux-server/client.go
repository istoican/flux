package main

import (
	"encoding/json"
	"flag"
	"io"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/istoican/flux"
)

const (
	defaultMaxMemory = 32 << 20
)

func main() {
	flag.Parse()

	path := flag.String("path", "./db", "database file")
	addr := flag.String("addr", "127.0.0.1:4000", "host")

	http.Handle("/", newServer(*path))

	log.Println("starting server")
	if err := http.ListenAndServe(*addr, nil); err != nil {
		panic(err)
	}
}

type server struct {
	*flux.DB
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Path[1:]

	w.Header().Set("Content-Type", "application/json")
	//w.(http.Flusher).Flush()

	switch r.Method {
	case "GET":
		if r.FormValue("watch") == "" {
			value, err := s.Get(key)
			response(w, value, err)
			return
		}

		watcher := s.Watch(key)

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
				response(w, e.Value, nil)
				w.(http.Flusher).Flush()
			}
		}
	case "POST":
		body, _ := ioutil.ReadAll(r.Body)
		//log.Println("VALUE: ", string(body))
		//r.ParseMultipartForm(defaultMaxMemory)
		//v := r.PostForm
		v := string(body)
		log.Println("VALUE: ", v)
		err := s.Put(key, v)
		response(w, v, err)
	case "DELETE":
		err := s.Delete(key)
		response(w, nil, err)
	}
}

func response(w io.Writer, value interface{}, err error) {
	if err := json.NewEncoder(w).Encode(value); err != nil {
		log.Printf("error %v \n", err)
	}
}

func newServer(path string) *server {
	db, err := flux.Open(path)
	if err != nil {
		panic(err)
	}
	return &server{db}
}

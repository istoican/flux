package http

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/istoican/flux"
)

func init() {
	handler := Handler{}

	http.Handle("/", handler)
}

// Handler :
type Handler struct {
	flux.Node
}

func (handler *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Path[1:]

	w.Header().Set("Content-Type", "application/json")
	//w.(http.Flusher).Flush()

	switch r.Method {
	case "GET":
		if r.FormValue("watch") == "" {
			value, err := server.Get(key)
			response(w, value, err)
			return
		}

		watcher := server.Watch(key)

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
		err := server.Put(key, v)
		response(w, v, err)
	case "DELETE":
		err := server.Delete(key)
		response(w, nil, err)
	}
}

func response(w io.Writer, value interface{}, err error) {
	if err := json.NewEncoder(w).Encode(value); err != nil {
		log.Printf("error %v \n", err)
	}
}

package main

import (
	"flag"
	"log"
	"net/http"
)

var server string

func init() {
	flag.StringVar(&server, "server", "", "server")
}

func main() {
	flag.Parse()

	http.HandleFunc("/", indexHandler)

	vars, err := Expvar(server)
	if err != nil {
		log.Panic(err)
	}
	stats.Nodes = vars.Flux.Members
	go record()

	log.Println("starting client")
	if err := http.ListenAndServe(":80", nil); err != nil {
		panic(err)
	}
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	tmpl.Execute(w, stats)
}

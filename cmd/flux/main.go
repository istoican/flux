package main

import (
	"flag"
	"log"
	"net/http"
	"strings"
)

var servers string

func init() {
	flag.StringVar(&servers, "servers", "", "servers")
}

func main() {
	flag.Parse()

	http.HandleFunc("/", indexHandler)

	go record(strings.Split(servers, ","))

	log.Println("starting client")
	if err := http.ListenAndServe(":80", nil); err != nil {
		panic(err)
	}
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	tmpl.Execute(w, stats)
}

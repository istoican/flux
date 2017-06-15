package main

import (
	"flag"
	"log"
	"net/http"
	_ "net/http/pprof"

	"github.com/istoican/flux"
	_ "github.com/istoican/flux/http"
)

var joinAddress string

func init() {
	flag.StringVar(&joinAddress, "join", "", "join address")
}

func main() {
	flag.Parse()

	if err := flux.Join(joinAddress); err != nil {
		log.Println(err)
	}

	log.Println("starting node")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}

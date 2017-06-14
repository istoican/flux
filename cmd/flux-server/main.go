package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/istoican/flux"
	"github.com/istoican/flux/store/memory"
)

const (
	defaultMaxMemory = 32 << 20
)

func main() {
	flag.Parse()

	address := flag.String("join", "", "join address")

	config := flux.Config{memory.NewStore()}
	server := flux.Start(config)

	if address != "" {
		server.join(address)
	}
	http.Handle("/", server)

	log.Println("starting server")
	if err := http.ListenAndServe(*addr, nil); err != nil {
		panic(err)
	}
}

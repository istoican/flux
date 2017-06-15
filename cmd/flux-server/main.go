package main

import (
	"flag"
	"log"
	"net/http"

	_ "github.com/istoican/flux/http"
)

func main() {
	flag.Parse()

	_ = flag.String("join", "", "join address")

	//if address != "" {
	//	server.join(address)
	//}

	log.Println("starting server")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}

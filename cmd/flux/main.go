package main

import (
	"flag"
	"log"
	"net/http"
)

var joinAddress string

func init() {
	flag.StringVar(&joinAddress, "join", "", "join address")
}

func main() {
	flag.Parse()

	http.HandleFunc("/", index)
	http.HandleFunc("/flux.json", flux)

	//f err := flux.Join(joinAddress); err != nil {
	//	log.Println(err)
	//}

	log.Println("starting client")
	if err := http.ListenAndServe(":80", nil); err != nil {
		panic(err)
	}
}

func index(w http.ResponseWriter, r *http.Request) {
	tmpl.Execute(w, nil)
}

func flux(w http.ResponseWriter, r *http.Request) {

}

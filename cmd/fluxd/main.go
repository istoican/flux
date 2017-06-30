package main

import (
	"flag"
	"log"
	"net/http"
	_ "net/http/pprof"

	"github.com/istoican/flux"
	"github.com/istoican/flux/consistent/hash"
	"github.com/istoican/flux/storage/disk"
	"github.com/istoican/flux/transport/http/handler"
)

var (
	joinAddress string
	storageType string
	hashFn      string
)

func init() {
	flag.StringVar(&joinAddress, "join", "", "join address")
	flag.StringVar(&storageType, "storage", "", "storage driver")
	flag.StringVar(&hashFn, "hashfn", "", "consistent hashing function")
}

func main() {
	flag.Parse()

	config := flux.DefaultConfig()

	if storageType == "disk" {
		db, err := disk.NewStore("flux.data")
		if err != nil {
			panic(err)
		}
		config.Store = db
	}
	if hashFn == "md5" {
		config.HashFn = hash.MD5
	}
	node, err := flux.New(config)
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		if err := node.Join(joinAddress); err != nil {
			log.Println(err)
		}
	}()

	http.Handle("/", handler.New(node))

	log.Println("starting node")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}

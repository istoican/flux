package main

import (
	"flag"
	"log"
	"net/http"
	_ "net/http/pprof"

	"github.com/istoican/flux"
	"github.com/istoican/flux/consistent/hash"
	"github.com/istoican/flux/storage/disk"
	transport "github.com/istoican/flux/transport/http"
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
	config.OnJoin = transport.Join
	config.OnLeave = transport.Leave
	config.Picker = transport.Handler()

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
	flux.Start(config)

	if err := flux.Join(joinAddress); err != nil {
		log.Println(err)
	}

	log.Println("starting node")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}

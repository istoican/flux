package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"runtime"
	"time"

	"github.com/istoican/flux"
)

var (
	stats []Stats
)

// Stats :
type Stats struct {
	Date  time.Time  `json:"date"`
	Stats []nodeInfo `json:"stats"`
}

type nodeInfo struct {
	Node      string `json:"node"`
	Memory    uint64 `json:"memory"`
	Reads     int64  `json:"reads"`
	Keys      int64  `json:"keys"`
	Deletions int64  `json:"deletions"`
	Inserts   int64  `json:"inserts"`
}

type expvars struct {
	Memstats runtime.MemStats
	Flux     flux.Stats
}

func record(servers []string) {
	for {
		var i []nodeInfo
		for _, server := range servers {
			info, err := read(server)
			if err != nil {
				log.Println(err)
			}
			i = append(i, info)
		}
		stats = append(stats, Stats{Date: time.Now(), Stats: i})
		//v, _ := json.Marshal(i)
		//log.Println("STATS: ", string(v))
		time.Sleep(1 * time.Second)
	}
}

func read(node string) (nodeInfo, error) {
	info := nodeInfo{Node: node}

	exp := expvars{}
	resp, err := http.Get("http://" + node + ":8080/debug/vars")
	if err != nil {
		return info, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	//log.Println("BODY: ", string(body))
	if err := json.Unmarshal(body, &exp); err != nil {
		return info, err
	}
	info.Memory = exp.Memstats.Alloc
	info.Reads = exp.Flux.Reads.Get()
	info.Deletions = exp.Flux.Deletions.Get()
	info.Inserts = exp.Flux.Inserts.Get()
	info.Keys = exp.Flux.Keys.Get()
	return info, nil
}
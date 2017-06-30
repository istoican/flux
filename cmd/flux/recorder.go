package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"reflect"
	"time"
)

var (
	stats []Stats
)

type Nodes map[string]string

type Timeline []time.Time

// Stats :
type Stats struct {
	Nodes    Nodes
	Timeline Timeline
	Metrics  map[string]Metrics
}

type Metrics struct {
	Memory    uint64 `json:"memory"`
	Reads     int64  `json:"reads"`
	Keys      int64  `json:"keys"`
	Deletions int64  `json:"deletions"`
	Inserts   int64  `json:"inserts"`
}

func record(server string) {
	if _, err := read(server); err != nil {
		log.Println(err)
	}
}

func read(node string) (nodeInfo, error) {
	info := nodeInfo{}
	vr := map[string]interface{}{}
	resp, err := http.Get("http://" + node + ":8080/debug/vars")
	if err != nil {
		return info, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err := json.Unmarshal(body, &vr); err != nil {
		return info, err
	}
	v := reflect.ValueOf(vr)

	vars := v.MapIndex(reflect.ValueOf("flux"))

	members := vars.Elem().MapIndex(reflect.ValueOf("members")).Elem()
	m := make([]string, 0)

	for _, k := range members.MapKeys() {
		v := members.MapIndex(k).Elem()
		m = append(m, k.String())
		m = append(m, v.String())
	}
	stats := vars.Elem().MapIndex(reflect.ValueOf("stats")).Elem()
	m2 := make([]string, 0)

	for _, k := range stats.MapKeys() {
		v := stats.MapIndex(k).Elem()
		m2 = append(m2, k.String())
		m2 = append(m2, v.String())
	}
	log.Println(m2)
	return info, nil
}

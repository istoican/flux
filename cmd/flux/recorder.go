package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"
)

var (
	stats  *Stats
	client http.Client
)

func init() {
	client = http.Client{
		Timeout: time.Duration(100 * time.Millisecond),
	}
	stats = &Stats{
		Nodes:   make(map[string]string),
		Metrics: make(map[time.Time]map[string]*Metrics),
	}
}

type expvar struct {
	Flux struct {
		Members map[string]string
		Stats   struct {
			Reads     uint64
			Keys      uint64
			Deletions uint64
			Inserts   uint64
		}
	}
	Memory struct {
		Alloc uint64
	} `json:"memstats"`
}

// Expvar :
func Expvar(node string) (expvar, error) {
	vars := expvar{}
	resp, err := http.Get("http://" + node + ":8080/debug/vars")
	if err != nil {
		return vars, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err := json.Unmarshal(body, &vars); err != nil {
		return vars, err
	}
	return vars, nil
}

// Stats :
type Stats struct {
	sync.RWMutex
	Nodes   map[string]string                 `json:"nodes"`
	Metrics map[time.Time]map[string]*Metrics `json:"metrics"`
}

// Metrics :
type Metrics struct {
	Memory    uint64 `json:"memory"`
	Reads     uint64 `json:"reads"`
	Keys      uint64 `json:"keys"`
	Deletions uint64 `json:"deletions"`
	Inserts   uint64 `json:"inserts"`
}

func (s *Stats) addNodes(nodes map[string]string) {
	s.Lock()
	defer s.Unlock()

	for k, v := range nodes {
		s.Nodes[k] = v
	}
}

func (s *Stats) addMetrics(date time.Time, metrics map[string]*Metrics) {
	s.Lock()
	defer s.Unlock()

	s.Metrics[date] = metrics
}

func record() {
	for {
		t := time.Now()
		metrics := make(map[string]*Metrics)
		var wg sync.WaitGroup
		for k, v := range stats.Nodes {
			wg.Add(1)
			go func() {
				defer wg.Done()
				vars, err := Expvar(v)
				if err != nil {
					log.Println(err)
					return
				}
				metrics[k] = &Metrics{
					Memory:    vars.Memory.Alloc,
					Inserts:   vars.Flux.Stats.Inserts,
					Keys:      vars.Flux.Stats.Keys,
					Reads:     vars.Flux.Stats.Reads,
					Deletions: vars.Flux.Stats.Deletions,
				}
				stats.addNodes(vars.Flux.Members)
			}()
		}
		wg.Wait()
		stats.addMetrics(t, metrics)

		time.Sleep(1 * time.Second)
	}
}

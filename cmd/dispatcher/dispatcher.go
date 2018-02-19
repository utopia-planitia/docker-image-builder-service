package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
	"sync"
)

type dispatcher struct {
	builders []*builder
	mutex    *sync.Mutex
	cond     *sync.Cond
}

func newDispatcher(endpoints []*url.URL, cpu, memory *int64, addr *string) *http.Server {

	builders := make([]*builder, len(endpoints))

	for i, e := range endpoints {
		r := httputil.NewSingleHostReverseProxy(e)
		builders[i] = &builder{
			name:     e.String(),
			proxy:    r,
			cpuquota: strconv.FormatInt(*cpu, 10),
			memory:   strconv.FormatInt(*memory, 10),
		}
	}

	m := &sync.Mutex{}
	c := sync.NewCond(m)

	s := &dispatcher{
		builders: builders,
		mutex:    m,
		cond:     c,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", s.handle)

	return &http.Server{
		Addr:    *addr,
		Handler: mux,
	}
}

func (s *dispatcher) handle(w http.ResponseWriter, r *http.Request) {

	log.Printf("requested path: %s\n", r.URL.Path)

	ip, err := parseClientIP(r)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		log.Printf("failed to parse client IP: %s\n", err)
		return
	}
	c := clientID(ip)

	b := s.selectWorker(c)
	defer s.recycle(b)
	b.handle(w, r)
}

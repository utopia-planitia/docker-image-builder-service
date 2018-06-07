package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
	"time"

	"github.com/damoon/ttlcache"
)

type dispatcher struct {
	builders []*builder
	cache    *cache
}

func newDispatcher(endpoints []*url.URL, cpu, memory *int64, network, addr *string) *http.Server {

	builders := make([]*builder, len(endpoints))

	for i, e := range endpoints {
		r := httputil.NewSingleHostReverseProxy(e)
		r.FlushInterval = 100 * time.Millisecond
		builders[i] = &builder{
			name:     e,
			proxy:    r,
			cpuquota: strconv.FormatInt(*cpu, 10),
			memory:   strconv.FormatInt(*memory, 10),
			network:  *network,
		}
	}

	cache := &cache{ttlcache.NewCache(time.Minute)}
	cache.StartCleanupTimer(10 * time.Second)

	s := &dispatcher{
		builders: builders,
		cache:    cache,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", s.handle)
	mux.HandleFunc("/healthz", healthz)

	return &http.Server{
		Addr:         *addr,
		Handler:      mux,
		ReadTimeout:  1800 * time.Second,
		WriteTimeout: 1800 * time.Second,
	}
}

func (s *dispatcher) handle(w http.ResponseWriter, r *http.Request) {

	log.Printf("requested path: %s\n", r.URL)

	ip, err := parseClientIP(r)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		log.Printf("failed to parse client IP: %s\n", err)
		return
	}
	c := clientID(ip)

	v := url.Values{}
	v["t"] = r.URL.Query()["t"]
	h := http.Header{}
	h["t"] = r.Header["GitBranchName"]

	b, err := s.selectWorker(r.Context(), c, v, h)
	if err != nil {
		log.Printf("failed select Worker: %s\n", err)
		return
	}
	defer s.returnWorker(b)
	b.handle(w, r)
}

package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

type tag string

type builder struct {
	docker *httputil.ReverseProxy
}

func newBuilder(endpoint *url.URL, addr *string) *http.Server {

	d := httputil.NewSingleHostReverseProxy(endpoint)
	b := &builder{
		docker: d,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", b.handle)

	return &http.Server{
		Addr:    *addr,
		Handler: mux,
	}
}

func (b *builder) handle(w http.ResponseWriter, r *http.Request) {

	log.Printf("requested path: %s\n", r.URL.Path)

	var t tag
	ts, ok := r.URL.Query()["t"]
	if ok {
		t = tag(ts[0])
	}
	log.Printf("building tag: %s\n", t)

	b.docker.ServeHTTP(w, r)
}

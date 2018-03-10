package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"regexp"
	"time"
)

var buildPath *regexp.Regexp

func init() {
	b, err := regexp.Compile("^/[^/]*/build")
	if err != nil {
		log.Fatalf("failed to prepare pattern matching: %s\n", err)
	}
	buildPath = b
}

type filename string

type builder struct {
	docker *httputil.ReverseProxy
}

func newBuilder(endpoint *url.URL, addr *string) *http.Server {

	d := httputil.NewSingleHostReverseProxy(endpoint)
	d.FlushInterval = 100 * time.Millisecond
	b := &builder{
		docker: d,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", b.handle)
	mux.HandleFunc("/uncachedBytes", uncachedBytes)

	return &http.Server{
		Addr:         *addr,
		Handler:      mux,
		ReadTimeout:  1800 * time.Second,
		WriteTimeout: 1800 * time.Second,
	}
}

func (b *builder) handle(w http.ResponseWriter, r *http.Request) {

	log.Printf("requested path: %s\n", r.URL)

	if buildPath.MatchString(r.URL.Path) {
		b.build(w, r)
	}

	b.docker.ServeHTTP(w, r)
}

func (b *builder) build(w http.ResponseWriter, r *http.Request) {
	tags, cacheFromBranches, currentBranch, err := parseTagsAndBranches(r)
	if err != nil {
		log.Printf("failed to parse tags and branches from request: %s\n", err)
	}

	values := r.URL.Query()
	values.Del("cachefrom")
	values.Set("pull", "1")
	values.Set("rm", "0")
	r.URL.RawQuery = values.Encode()

	load(tags, cacheFromBranches, currentBranch)
	b.docker.ServeHTTP(w, r)
	save(tags, currentBranch)
}

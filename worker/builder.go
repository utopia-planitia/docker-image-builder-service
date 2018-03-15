package main

import (
	"encoding/json"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"regexp"
	"time"
)

var (
	buildPath     *regexp.Regexp
	containerPath *regexp.Regexp
)

func init() {
	b, err := regexp.Compile("^/([^/]*/)?build")
	if err != nil {
		log.Fatalf("failed to prepare pattern matching: %s\n", err)
	}
	buildPath = b

	c, err := regexp.Compile("^/([^/]*/)?containers")
	if err != nil {
		log.Fatalf("failed to prepare pattern matching: %s\n", err)
	}
	containerPath = c
}

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

	if isRequestingContainer(r.URL.Path) {
		w.WriteHeader(http.StatusForbidden)
		log.Printf("container use is forbidden: %s\n", r.URL)
		return
	}

	if isRequestingBuild(r.URL.Path) {
		b.build(w, r)
		return
	}

	b.docker.ServeHTTP(w, r)
}

func isRequestingContainer(p string) bool {
	return containerPath.MatchString(p)
}

func isRequestingBuild(p string) bool {
	return buildPath.MatchString(p)
}

func (b *builder) build(w http.ResponseWriter, r *http.Request) {

	tags, cacheFromBranches, currentBranch, err := parseTagsAndBranches(r)
	if err != nil {
		log.Printf("failed to parse tags and branches from request: %s\n", err)
	}
	if len(tags) == 0 {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("untagged builds are not supported"))
		log.Printf("no tags set: %s\n", r.URL)
		return
	}

	load(tags, cacheFromBranches, currentBranch)
	cacheFromLocalImages(r, cacheSources(tags, cacheFromBranches, currentBranch))
	log.Printf("docker forwarded request: %v\n", r)
	b.docker.ServeHTTP(w, r)
	save(tags, currentBranch)
}

func cacheFromLocalImages(r *http.Request, tags []*tag) {

	values := r.URL.Query()

	values.Set("pull", "1")
	values.Set("rm", "0")

	values.Del("cachefrom")

	possibleCaches := []string{}
	for _, tag := range tags {
		possibleCaches = append(possibleCaches, "cache:5000/"+tag.String())
	}
	cachefrom, err := json.Marshal(possibleCaches)
	if err != nil {
		log.Printf("failed to marshal local images: %s\n", err)
	}
	if cachefrom != nil {
		values.Set("cachefrom", string(cachefrom))
		log.Printf("added localimage to cachefrom: %s\n", string(cachefrom))
	}

	r.URL.RawQuery = values.Encode()
}

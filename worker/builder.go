package main

import (
	"encoding/json"
	"fmt"
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
	tagPath       *regexp.Regexp
)

func init() {
	b, err := regexp.Compile("^/([^/]*/)?build")
	if err != nil {
		log.Fatalf("failed to prepare build pattern: %s\n", err)
	}
	buildPath = b

	c, err := regexp.Compile("^/([^/]*/)?containers")
	if err != nil {
		log.Fatalf("failed to prepare containers pattern: %s\n", err)
	}
	containerPath = c

	t, err := regexp.Compile("^/([^/]*/)?images/([^/]*)/tag")
	if err != nil {
		log.Fatalf("failed to prepare tag pattern: %s\n", err)
	}
	tagPath = t
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
	mux.HandleFunc("/uncachedSize", uncachedSize)
	mux.HandleFunc("/healthz", healthz)

	return &http.Server{
		Addr:         *addr,
		Handler:      mux,
		ReadTimeout:  1800 * time.Second,
		WriteTimeout: 1800 * time.Second,
	}
}

func (b *builder) handle(w http.ResponseWriter, r *http.Request) {

	log.Printf("requested path: %s\n", r.URL)

	if isRequestingContainer(r.URL.Path) && r.Method != "GET" {
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

func isRequestingContainer(r string) bool {
	return containerPath.MatchString(r)
}

func isRequestingBuild(r string) bool {
	return buildPath.MatchString(r)
}

func (b *builder) build(w http.ResponseWriter, r *http.Request) {

	tag, currentBranch, headBranch, err := parseTagsAndBranches(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		message := fmt.Sprintf("failed to parse request: %s\n", err)
		_, err := w.Write([]byte(message))
		if err != nil {
			log.Printf("failed to write to response: %v\n", err)
		}
		return
	}

	cachedSources := restoreCache(tag, currentBranch, headBranch)
	configureRequest(r, cachedSources)
	log.Printf("docker forwarded request: %v\n", r)
	b.docker.ServeHTTP(w, r)
	archiveCache(tag, currentBranch)
}

func configureRequest(r *http.Request, tags []*tag) {

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
		log.Printf("failed to marshal cachefrom: %s\n", err)
	}
	if cachefrom != nil {
		values.Set("cachefrom", string(cachefrom))
		log.Printf("added localimage to cachefrom: %s\n", string(cachefrom))
	}

	r.URL.RawQuery = values.Encode()
}

package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
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

	if isRequestingTag(r.URL.Path) {
		b.tag(w, r)
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

func isRequestingTag(r string) bool {
	return tagPath.MatchString(r)
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
	cacheFromLocalImages(r)
	log.Printf("docker forwarded request: %v\n", r)
	b.docker.ServeHTTP(w, r)
	save(tags, currentBranch)
}

func (b *builder) tag(w http.ResponseWriter, r *http.Request) {
	repo := r.URL.Query().Get("repo")
	t := r.URL.Query().Get("tag")
	tags := []*tag{&tag{image: repo, version: t}}
	b.docker.ServeHTTP(w, r)
	save(tags, "")
}

func cacheFromLocalImages(r *http.Request) {

	values := r.URL.Query()

	values.Set("pull", "1")
	values.Set("rm", "0")

	values.Del("cachefrom")
	cachefrom, err := json.Marshal(localImages())
	if err != nil {
		log.Printf("failed to marshal local images: %s\n", err)
	}
	if cachefrom != nil {
		values.Set("cachefrom", string(cachefrom))
		log.Printf("added localimage to cachefrom: %s\n", string(cachefrom))
	}

	r.URL.RawQuery = values.Encode()
}

func localImages() []string {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		log.Printf("failed to pull image: %s", err)
		return nil
	}

	images, err := cli.ImageList(context.Background(), types.ImageListOptions{})
	if err != nil {
		log.Printf("failed list images: %s", err)
	}
	tags := []string{}
	for _, image := range images {
		for _, name := range image.RepoTags {
			if strings.HasPrefix(name, "cache:5000/") {
				tags = append(tags, name)
			}
		}
	}
	return tags
}

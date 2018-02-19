package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os/exec"
	"strings"

	"github.com/damoon/docker-image-builder-service/dibs"
)

type filename string

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

	t, err := dibs.ParseTag(r)
	if err != nil {
		log.Printf("cache preparation failed: %s\n", err)
		b.docker.ServeHTTP(w, r)
		return
	}

	log.Printf("building tag: %s\n", t)

	f := cachedFilename(t)

	load(f)

	b.docker.ServeHTTP(w, r)

	if t.Version == "latest" {
		save(t, f)
	}
}

func load(f filename) {
	log.Printf("loading cached file %s", string(f))
	output, err := exec.Command("load", string(f)).CombinedOutput()
	if err != nil {
		log.Printf("loading cached file %s failed: %v: %v", string(f), err, string(output))
	}
}

func save(t *dibs.Tag, f filename) {
	log.Printf("saving image %s to file %s", t.String(), string(f))
	output, err := exec.Command("save", t.String(), string(f)).CombinedOutput()
	if err != nil {
		log.Printf("saving image %s to file %s failed: %v: %v", t.String(), string(f), err, string(output))
	}
}

func cachedFilename(t *dibs.Tag) filename {
	return filename(strings.Replace(t.Image, "/", "~", -1) + ":" + t.Version)
}

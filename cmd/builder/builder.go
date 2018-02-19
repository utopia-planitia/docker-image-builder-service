package main

import (
	"encoding/json"
	"errors"
	"fmt"
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

	log.Printf("url: %s\n", r.URL.String())

	log.Printf("requested path: %s\n", r.URL.Path)

	t, err := dibs.ParseTag(r)
	if err != nil {
		log.Printf("tag cache preparation failed: %s\n", err)
	}
	var f filename
	if t != nil {
		log.Printf("building tag: %s\n", t)
		f = cachedLatestFilename(t)
		load(f)
	}

	cf, err := parseCachefrom(r)
	if err != nil {
		log.Printf("cachefrom preparation failed: %s\n", err)
	}
	if cf != nil {
		log.Printf("cachefrom: %s\n", cf)
		for _, e := range cf {
			load(cachedBranchFilename(t, e))
		}
	}
	values := r.URL.Query()
	values.Del("cachefrom")
	r.URL.RawQuery = values.Encode()

	b.docker.ServeHTTP(w, r)

	if t != nil && t.Version == "latest" {
		save(t, f)
	}
	if cf != nil {
		for _, e := range cf {
			save(t, cachedBranchFilename(t, e))
		}
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

func cachedLatestFilename(t *dibs.Tag) filename {
	return filename(strings.Replace(t.Image, "/", "~", -1) + ":latest")
}

func cachedBranchFilename(t *dibs.Tag, bn string) filename {
	return filename(strings.Replace(t.Image, "/", "~", -1) + ":" + bn)
}

func parseCachefrom(r *http.Request) ([]string, error) {
	cf, ok := r.URL.Query()["cachefrom"]
	if !ok {
		return nil, errors.New("parameter cachefrom not set")
	}
	if len(cf) != 1 {
		return nil, errors.New("parameter cachefrom not set exactly once")
	}

	sr := strings.NewReader(cf[0])
	j := json.NewDecoder(sr)
	var l []string
	err := j.Decode(&l)
	if err != nil {
		return nil, fmt.Errorf("failed to json decode cachefrom: %s", err)
	}
	return l, nil
}

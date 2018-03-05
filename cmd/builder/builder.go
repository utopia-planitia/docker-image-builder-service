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
	"regexp"
	"strings"
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

	if !buildPath.MatchString(r.URL.Path) {
		b.docker.ServeHTTP(w, r)
		return
	}

	t, err := parseTag(r)
	if err != nil {
		log.Printf("tag cache preparation failed: %s\n", err)
		b.docker.ServeHTTP(w, r)
		return
	}

	log.Printf("building tag: %s\n", t)
	f := cachedLatestFilename(t)
	cf, err := parseCachefromBranches(r)
	if err != nil {
		log.Printf("cachefrom preparation failed: %s\n", err)
	}
	if len(cf) != 0 {
		values := r.URL.Query()
		values.Del("cachefrom")
		r.URL.RawQuery = values.Encode()
	}

	load(t, cf, f)
	b.docker.ServeHTTP(w, r)
	save(t, cf, f)
}

func load(t *tag, cf []string, f filename) {
	loadCommand(t, f)
	if cf != nil {
		log.Printf("cachefrom: %s\n", cf)
		for _, e := range cf {
			if e == "master" {
				continue
			}
			loadCommand(t, cachedBranchFilename(t, e))
		}
	}
}

func save(t *tag, cf []string, f filename) {
	if t.version == "latest" {
		saveCommand(t, f)
	}
	for _, e := range cf {
		if e == "master" {
			continue
		}
		saveCommand(t, cachedBranchFilename(t, e))
	}
}

func loadCommand(t fmt.Stringer, f filename) {
	log.Printf("loading cached file %s", f)
	/* #nosec */
	output, err := exec.Command("load", t.String(), string(f)).CombinedOutput()
	if err != nil {
		log.Printf("loading cached file %s failed: %v: %v", f, err, string(output))
	}
}

func saveCommand(t fmt.Stringer, f filename) {
	log.Printf("saving image %s to file %s", t, f)
	/* #nosec */
	output, err := exec.Command("save", t.String(), string(f)).CombinedOutput()
	if err != nil {
		log.Printf("saving image %s to file %s failed: %v: %v", t, f, err, string(output))
	}
}

func cachedLatestFilename(t *tag) filename {
	return filename(strings.Replace(t.image, "/", "~", -1) + ":latest")
}

func cachedBranchFilename(t *tag, bn string) filename {
	return filename(strings.Replace(t.image, "/", "~", -1) + ":" + bn)
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

func parseCachefromBranches(r *http.Request) ([]string, error) {
	cf, err := parseCachefrom(r)
	if err != nil {
		return nil, fmt.Errorf("failed to parse branches: %s", err)
	}
	var l []string
	for _, e := range cf {
		if strings.HasPrefix(e, "branch=") {
			l = append(l, e[len("branch="):])
		}
	}
	return l, nil
}

package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"regexp"
)

var buildPath *regexp.Regexp

func init() {
	b, err := regexp.Compile("^/[^/]*/build")
	if err != nil {
		log.Fatalf("failed to prepare pattern matching: %s\n", err)
	}
	buildPath = b
}

type builder struct {
	name            *url.URL
	proxy           *httputil.ReverseProxy
	cpuquota        string
	memory          string
	openConnections int32
	dedicatedTo     clientID
	lastestUse      int64
}

func (b *builder) handle(w http.ResponseWriter, r *http.Request) {
	if buildPath.MatchString(r.URL.Path) {
		b.configureBuildRequest(r)
	}
	b.proxy.ServeHTTP(w, r)
}

func (b *builder) configureBuildRequest(r *http.Request) {
	values := r.URL.Query()
	values.Set("cpuquota", b.cpuquota)
	values.Set("memory", b.memory)
	values.Set("networkmode", "host")
	r.URL.RawQuery = values.Encode()
}

package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"regexp"

	"github.com/damoon/docker-image-builder-service/dibs"
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
	name            string
	proxy           *httputil.ReverseProxy
	cpuquota        string
	memory          string
	openConnections int32
	dedicatedTo     clientID
	lastestUse      int64
}

func (b *builder) handle(t *dibs.Tag, w http.ResponseWriter, r *http.Request) {

	if buildPath.MatchString(r.URL.Path) {
		values := r.URL.Query()
		values.Set("cpuquota", b.cpuquota)
		values.Set("memory", b.memory)
		r.URL.RawQuery = values.Encode()
	}

	b.proxy.ServeHTTP(w, r)
}

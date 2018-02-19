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
	buildResources  string
	openConnections int32
	dedicatedTo     clientID
	lastestUse      int64
}

func (b *builder) handle(t *dibs.Tag, w http.ResponseWriter, r *http.Request) {

	if buildPath.MatchString(r.URL.Path) {
		r.URL.RawQuery += b.buildResources
	}

	b.proxy.ServeHTTP(w, r)
}

package main

import (
	"fmt"
	"log"
	"net"
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

type builder struct {
	name            *url.URL
	proxy           *httputil.ReverseProxy
	cpuquota        string
	memory          string
	network         string
	openConnections int32
	dedicatedTo     clientID
	lastestUse      int64
}

func (b *builder) healthy() bool {
	var netClient = &http.Client{
		Timeout: time.Second * 10,
		Transport: &http.Transport{
			Dial: (&net.Dialer{
				Timeout: 5 * time.Second,
			}).Dial,
			TLSHandshakeTimeout: 5 * time.Second,
		},
	}

	response, err := netClient.Get(fmt.Sprintf("%s/healthy", b.name))
	if err != nil {
		log.Printf("health check %s: %v", b.name, err)
		return false
	}

	defer response.Body.Close()

	if response.StatusCode == http.StatusOK {
		log.Printf("health check %s: status code %d", b.name, response.StatusCode)
		return false
	}

	return true
}

func (b *builder) handle(w http.ResponseWriter, r *http.Request) {
	if buildPath.MatchString(r.URL.Path) {
		b.configureBuildRequest(r)
	}
	b.proxy.ServeHTTP(w, r)
}

func (b *builder) configureBuildRequest(r *http.Request) {
	values := r.URL.Query()
	if b.cpuquota != "0" {
		values.Set("cpuquota", b.cpuquota)
	}
	if b.memory != "0" {
		values.Set("memory", b.memory)
	}
	if b.network != "" {
		values.Set("networkmode", b.network)
	}
	r.URL.RawQuery = values.Encode()
}

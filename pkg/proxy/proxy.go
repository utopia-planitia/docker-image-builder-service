package proxy

import (
	"context"
	"log"
	"net/http"
	"net/http/httputil"
	"regexp"
	"strconv"

	"golang.org/x/sync/semaphore"
)

// QueuedProxy structures the data to proxy requests
type QueuedProxy struct {
	proxy            *httputil.ReverseProxy
	queue            *semaphore.Weighted
	buildResources   string
	queuedURLPattern *regexp.Regexp
}

// New creates new instances of the QueneProxy handler
func New(target *httputil.ReverseProxy, parallelism, cpu, memory int64, pattern *regexp.Regexp) *QueuedProxy {

	return &QueuedProxy{
		proxy:            target,
		queue:            semaphore.NewWeighted(parallelism),
		buildResources:   "&cpuquota=" + strconv.FormatInt(cpu, 10) + "&memory=" + strconv.FormatInt(memory, 10),
		queuedURLPattern: pattern,
	}
}

// Handle processes http requests
func (p *QueuedProxy) Handle(w http.ResponseWriter, r *http.Request) {

	log.Printf("requested path: %s\n", r.URL.Path)

	// skip none build requests
	if !p.queuedURLPattern.MatchString(r.URL.Path) {
		log.Printf("skip queue for: %s\n", r.URL.Path)
		p.proxy.ServeHTTP(w, r)
		return
	}

	// add resource limit to build
	r.URL.RawQuery += p.buildResources

	t, ok := r.URL.Query()["t"]
	image := "undefined"
	if ok && len(t) >= 1 {
		image = t[0]
	}

	//	// forward request with free build slot
	//	if p.queue.TryAcquire(1) {
	//		defer p.queue.Release(1)
	//		defer log.Printf("finished building image: %s\n", image)
	//		log.Printf("building image: %s\n", image)
	//		p.proxy.ServeHTTP(w, r)
	//		return
	//	}

	log.Printf("queued image: %s\n", image)
	p.queue.Acquire(context.Background(), 1)

	// serialize body until build happens?

	p.buildImage(w, r, image)
}

func (p *QueuedProxy) buildImage(w http.ResponseWriter, r *http.Request, image string) {
	defer p.queue.Release(1)
	defer log.Printf("finished building image: %s\n", image)
	log.Printf("building image: %s\n", image)
	p.proxy.ServeHTTP(w, r)
}

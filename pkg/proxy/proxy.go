package proxy

import (
	"log"
	"net/http/httputil"
	"net/http"
	"golang.org/x/sync/semaphore"
	"context"
	"time"
)

type QueuedProxy struct {
	proxy         *httputil.ReverseProxy
	queue         *semaphore.Weighted
}

func New(target *httputil.ReverseProxy, parallelism int64) *QueuedProxy {
	return &QueuedProxy{
		proxy:         target,
		queue:         semaphore.NewWeighted(parallelism),
	}
}

func (p *QueuedProxy) Handle(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("X-GoProxy", "GoProxy")
	r.Header.Set("a","b")
	r.URL.RawQuery += "&n=m"
	log.Printf("request: %s\n", r.URL)
	p.queue.Acquire(context.Background(), 1)
	time.Sleep(1000 * time.Millisecond)
	p.proxy.ServeHTTP(w, r)
	p.queue.Release(1)
}

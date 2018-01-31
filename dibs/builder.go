package dibs

import (
	"log"
	"net/http"
	"net/http/httputil"
)

type builder struct {
	proxy          *httputil.ReverseProxy
	buildResources string
}

func (b *builder) build(t tag, w http.ResponseWriter, r *http.Request) {

	// add resource limit to build
	r.URL.RawQuery += b.buildResources

	// build images in free slot
	log.Printf("building image: %s\n", t)
	b.proxy.ServeHTTP(w, r)
	log.Printf("finished building image: %s\n", t)

}

func (b *builder) forward(w http.ResponseWriter, r *http.Request) {
	b.proxy.ServeHTTP(w, r)
}

package dibs

import (
	"errors"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"regexp"
	"strconv"
	"sync"
)

var buildPath *regexp.Regexp

type scheduler struct {
	builders []*builder
	mutex    *sync.Mutex
}

func init() {
	b, err := regexp.Compile("^/[^/]*/build")
	if err != nil {
		log.Fatalf("failed to prepare pattern matching: %s\n", err)
	}
	buildPath = b
}

// NewScheduler creates a new image builds scheduling http server
func NewScheduler(endpoints []*url.URL, cpu, memory *int64, addr *string) *http.Server {

	builders := make([]*builder, len(endpoints))

	for i, e := range endpoints {
		r := httputil.NewSingleHostReverseProxy(e)
		builders[i] = &builder{
			proxy:          r,
			buildResources: "&cpuquota=" + strconv.FormatInt(*cpu, 10) + "&memory=" + strconv.FormatInt(*memory, 10),
		}
	}

	s := &scheduler{
		builders: builders,
		mutex:    &sync.Mutex{},
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", s.Handle)

	return &http.Server{
		Addr:    *addr,
		Handler: mux,
	}
}

// Handle processes http requests
func (s *scheduler) Handle(w http.ResponseWriter, r *http.Request) {

	log.Printf("requested path: %s\n", r.URL.Path)

	var t tag
	ts, ok := r.URL.Query()["t"]
	if ok {
		t = tag(ts[0])
	}

	ip, err := parseClientIP(r)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write([]byte(err.Error()))
		log.Printf("%s\n", err)
		return
	}

	log.Printf("queued image: %s\n", t)

	c := clientID(ip)
	b := s.selectBuilder(t, c)
	defer b.Close()

	if buildPath.MatchString(r.URL.Path) {
		b.build(t, w, r)
		return
	}

	b.forward(w, r)
}

func parseTag(v url.Values) (tag, error) {
	t, ok := v["t"]
	if !ok {
		return tag(""), errors.New("missing parameter t")
	}
	return tag(t[0]), nil
}

package main

import (
	"errors"
	"net/http"
	"strings"
)

type tag struct {
	image, version string
}

func (t tag) String() string {
	return t.image + ":" + t.version
}

func parseTag(r *http.Request) (*tag, error) {
	t, ok := r.URL.Query()["t"]
	if !ok {
		return &tag{}, errors.New("parameter t not set")
	}
	if len(t) != 1 {
		return &tag{}, errors.New("parameter t not set exactly once")
	}
	return newTag(t[0])
}

func newTag(t string) (*tag, error) {
	i := strings.LastIndex(t, ":")
	if i == -1 {
		return &tag{
			image:   t,
			version: "latest",
		}, nil
	}
	if i < strings.LastIndex(t, "/") {
		return &tag{
			image:   t,
			version: "latest",
		}, nil
	}
	return &tag{
		image:   t[:i],
		version: t[i+1:],
	}, nil
}

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
	s := strings.Split(t, ":")
	if len(s) == 2 {
		return &tag{
			image:   s[0],
			version: s[1],
		}, nil
	}
	if len(s) == 1 {
		return &tag{
			image:   s[0],
			version: "latest",
		}, nil
	}
	return &tag{}, errors.New("tag is malformed")
}

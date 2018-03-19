package main

import (
	"strings"
)

type tag struct {
	image, version string
}

func (t tag) String() string {
	return t.image + ":" + t.version
}

func newTag(t string) *tag {
	i := strings.LastIndex(t, ":")
	if i == -1 {
		return &tag{
			image:   t,
			version: "latest",
		}
	}
	if i < strings.LastIndex(t, "/") {
		return &tag{
			image:   t,
			version: "latest",
		}
	}
	return &tag{
		image:   t[:i],
		version: t[i+1:],
	}
}

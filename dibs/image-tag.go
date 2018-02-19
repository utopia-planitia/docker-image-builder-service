package dibs

import (
	"errors"
	"net/http"
	"strings"
)

// Tag represents a docker image name
type Tag struct {
	Image, Version string
}

func (t Tag) String() string {
	return t.Image + ":" + t.Version
}

// ParseTag parses a Build Tag from a request
func ParseTag(r *http.Request) (*Tag, error) {
	t, ok := r.URL.Query()["t"]
	if !ok {
		return &Tag{}, errors.New("parameter t not set")
	}
	if len(t) != 1 {
		return &Tag{}, errors.New("parameter t not set exactly once")
	}
	return newTag(t[0])
}

func newTag(t string) (*Tag, error) {
	s := strings.Split(t, ":")
	if len(s) == 2 {
		return &Tag{
			Image:   s[0],
			Version: s[1],
		}, nil
	}
	if len(s) == 1 {
		return &Tag{
			Image:   s[0],
			Version: "latest",
		}, nil
	}
	return &Tag{}, errors.New("tag is malformed")
}

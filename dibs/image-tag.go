package dibs

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
)

// Tag represents a docker image name
type Tag struct {
	Image, Version string
}

// ParseTag creates a Tag from a request
func ParseTag(r *http.Request) (*Tag, error) {
	s, err := stringFromRequest(r)
	if err != nil {
		return &Tag{}, fmt.Errorf("failed to parse tag: %s", err)
	}
	ss := strings.Split(s, ":")
	if len(ss) == 2 {
		return &Tag{
			Image:   ss[0],
			Version: ss[1],
		}, nil
	}
	if len(ss) == 1 {
		return &Tag{
			Image:   ss[0],
			Version: "latest",
		}, nil
	}
	return &Tag{}, errors.New("parameter t is malformed")
}

func stringFromRequest(r *http.Request) (string, error) {
	ts, ok := r.URL.Query()["t"]
	if !ok {
		return "", errors.New("parameter t not set")
	}
	return ts[0], nil
}

func (t Tag) String() string {
	return t.Image + ":" + t.Version
}

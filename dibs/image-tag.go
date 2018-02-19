package dibs

import (
	"encoding/json"
	"errors"
	"fmt"
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

// ParseCachefrom creates a list of cache from Tags from a request
func ParseCachefrom(r *http.Request) ([]string, error) {
	cf, ok := r.URL.Query()["cachefrom"]
	if !ok {
		return nil, errors.New("parameter cachefrom not set")
	}
	if len(cf) != 1 {
		return nil, errors.New("parameter cachefrom not set exactly once")
	}

	sr := strings.NewReader(cf[0])
	j := json.NewDecoder(sr)
	var l []string
	err := j.Decode(&l)
	if err != nil {
		return nil, fmt.Errorf("failed to json decode cachefrom: %s", err)
	}
	return l, nil
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

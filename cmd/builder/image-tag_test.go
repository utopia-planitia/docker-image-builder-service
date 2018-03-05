package main

import (
	"strings"
	"testing"
)

var tagTests = []struct {
	in  string
	out *tag
	err bool
}{
	{"cassandra", &tag{image: "cassandra", version: "latest"}, false},
	{"cassandra:latest", &tag{image: "cassandra", version: "latest"}, false},
	{"cassandra:2", &tag{image: "cassandra", version: "2"}, false},

	{"user/image:2", &tag{image: "user/image", version: "2"}, false},

	{"registry.domain.tld/cassandra:2", &tag{image: "registry.domain.tld/cassandra", version: "2"}, false},
	{"registry.domain.tld/user/image:2", &tag{image: "registry.domain.tld/user/image", version: "2"}, false},

	{"registry.domain.tld:5000/image:3", &tag{image: "registry.domain.tld:5000/image", version: "3"}, false},
	{"registry.domain.tld:5000/image:latest", &tag{image: "registry.domain.tld:5000/image", version: "latest"}, false},
	{"registry.domain.tld:5000/image", &tag{image: "registry.domain.tld:5000/image", version: "latest"}, false},

	{"registry.domain.tld:5000/user/image:3", &tag{image: "registry.domain.tld:5000/user/image", version: "3"}, false},
	{"registry.domain.tld:5000/user/image:latest", &tag{image: "registry.domain.tld:5000/user/image", version: "latest"}, false},
	{"registry.domain.tld:5000/user/image", &tag{image: "registry.domain.tld:5000/user/image", version: "latest"}, false},
}

func TestNewTag(t *testing.T) {
	for _, tt := range tagTests {
		tag, err := newTag(tt.in)
		if (err != nil) != tt.err {
			t.Errorf("newTag(%q) returned an error: %s", tt.in, err)
			continue
		}
		if !compare(tag, tt.out) {
			t.Errorf("newTag(%q) => %q, want %q", tt.in, debugPrint(tag), debugPrint(tt.out))
		}
	}
}

func debugPrint(t *tag) string {
	return "{image: " + t.image + ", version: " + t.version + "}"
}

func compare(a, b *tag) bool {
	if &a == &b {
		return true
	}
	if strings.Compare(a.image, b.image) != 0 {
		return false
	}
	if strings.Compare(a.version, b.version) != 0 {
		return false
	}
	return true
}

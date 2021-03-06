package main

import (
	"strings"
	"testing"
)

var tagTests = []struct {
	in  string
	out *tag
}{
	{"cassandra", &tag{image: "cassandra", version: "latest"}},
	{"cassandra:latest", &tag{image: "cassandra", version: "latest"}},
	{"cassandra:2", &tag{image: "cassandra", version: "2"}},

	{"user/image:2", &tag{image: "user/image", version: "2"}},

	{"registry.domain.tld/cassandra:2", &tag{image: "registry.domain.tld/cassandra", version: "2"}},
	{"registry.domain.tld/user/image:2", &tag{image: "registry.domain.tld/user/image", version: "2"}},

	{"registry.domain.tld:5000/image:3", &tag{image: "registry.domain.tld:5000/image", version: "3"}},
	{"registry.domain.tld:5000/image:latest", &tag{image: "registry.domain.tld:5000/image", version: "latest"}},
	{"registry.domain.tld:5000/image", &tag{image: "registry.domain.tld:5000/image", version: "latest"}},

	{"registry.domain.tld:5000/user/image:3", &tag{image: "registry.domain.tld:5000/user/image", version: "3"}},
	{"registry.domain.tld:5000/user/image:latest", &tag{image: "registry.domain.tld:5000/user/image", version: "latest"}},
	{"registry.domain.tld:5000/user/image", &tag{image: "registry.domain.tld:5000/user/image", version: "latest"}},
}

func TestNewTag(t *testing.T) {
	for _, tt := range tagTests {
		tag := newTag(tt.in)
		if !tag.equals(tt.out) {
			t.Errorf("newTag(%q) => %q, want %q", tt.in, tag.debugPrint(), tt.out.debugPrint())
		}
	}
}

func (t *tag) debugPrint() string {
	return "{image: " + t.image + ", version: " + t.version + "}"
}

func (t *tag) equals(b *tag) bool {
	if &t == &b {
		return true
	}
	if strings.Compare(t.image, b.image) != 0 {
		return false
	}
	if strings.Compare(t.version, b.version) != 0 {
		return false
	}
	return true
}

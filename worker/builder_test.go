package main

import (
	"testing"
)

var buildTests = []struct {
	in  string
	out bool
}{
	{"/", false},
	{"/_ping", false},
	{"/build", true},
	{"/v1.24/build", true},
	{"/containers", false},
	{"/v1.24/containers", false},
	{"/images/name/tag", false},
	{"/v1.24/images/name/tag", false},
}

var containerTests = []struct {
	in  string
	out bool
}{
	{"/", false},
	{"/_ping", false},
	{"/build", false},
	{"/v1.24/build", false},
	{"/containers", true},
	{"/v1.24/containers", true},
	{"/images/name/tag", false},
	{"/v1.24/images/name/tag", false},
}

func TestIsRequestingBuild(t *testing.T) {
	for _, tt := range buildTests {
		b := isRequestingBuild(tt.in)
		if b != tt.out {
			t.Errorf("isRequestingBuild(%q) => %t, want %t", tt.in, b, tt.out)
		}
	}
}

func TestIsRequestingContainer(t *testing.T) {
	for _, tt := range containerTests {
		b := isRequestingContainer(tt.in)
		if b != tt.out {
			t.Errorf("isRequestingBuild(%q) => %t, want %t", tt.in, b, tt.out)
		}
	}
}

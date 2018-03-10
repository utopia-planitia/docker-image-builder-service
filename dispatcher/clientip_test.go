package main

import (
	"net/http"
	"strings"
	"testing"
)

var tagTests = []struct {
	in  http.Request
	out string
}{
	{http.Request{Header: http.Header{"X-Forwarded-For": []string{"8.8.8.8"}}}, "8.8.8.8"},
	{http.Request{Header: http.Header{"X-Real-Ip": []string{"8.8.8.8"}}}, "8.8.8.8"},
	{http.Request{
		Header: http.Header{
			"X-Forwarded-For": []string{"8.8.8.8"},
			"X-Real-Ip":       []string{"8.8.4.4"},
		}},
		"8.8.8.8"},
	{http.Request{
		Header: http.Header{
			"X-Forwarded-For": []string{"8.8.8.8"},
			"X-Real-Ip":       []string{"8.8.4.4"},
		}},
		"8.8.8.8"},
	{http.Request{
		Header: http.Header{
			"X-Forwarded-For": []string{"8.8.8.8"},
			"X-Real-Ip":       []string{"8.8.4.4"},
		},
		RemoteAddr: "8.8.8.1"},
		"8.8.8.8"},
	{http.Request{
		Header: http.Header{
			"X-Real-Ip": []string{"8.8.4.4"},
		},
		RemoteAddr: "8.8.8.1"},
		"8.8.4.4"},
	{http.Request{
		RemoteAddr: "8.8.8.1:123"},
		"8.8.8.1"},
}

func TestParseClientIP(t *testing.T) {
	for _, tt := range tagTests {
		ip, err := parseClientIP(&tt.in)
		if err != nil {
			t.Errorf("parseClientIP(%v) errored: %v", tt.in, err)
			continue
		}
		if strings.Compare(ip, tt.out) != 0 {
			t.Errorf("parseClientIP(%v) => %q, want %q", tt.in, ip, tt.out)
		}
	}
}

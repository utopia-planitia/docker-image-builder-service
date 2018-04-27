package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

func parseTagsAndBranches(r *http.Request) (*tag, string, error) {

	ts := r.URL.Query()["t"]
	if len(ts) != 1 {
		return &tag{}, "", errors.New("tag parameter is not set exactly once")
	}

	if err := errorOnCacheFromIsSet(r); err != nil {
		return &tag{}, "", err
	}

	branch := r.Header.Get("GitBranchName")
	if branch == "" {
		return &tag{}, "", errors.New("Branch not set via http header 'GitBranchName'")
	}

	return newTag(ts[0]), branch, nil
}

func errorOnCacheFromIsSet(r *http.Request) error {
	cfs := r.URL.Query()["cachefrom"]
	if len(cfs) != 1 {
		return errors.New("cachefrom parameter not set exactly one")
	}
	cf, err := decodeCachefromJSON(cfs[0])
	if err != nil {
		return err
	}
	if len(cf) != 0 {
		return errors.New("the use of --cache-from is not supported")
	}
	return nil
}

func decodeCachefromJSON(cf string) ([]string, error) {
	r := strings.NewReader(cf)
	j := json.NewDecoder(r)
	var l []string
	err := j.Decode(&l)
	if err != nil {
		return nil, fmt.Errorf("failed to json decode cachefrom: %s", err)
	}
	return l, nil
}

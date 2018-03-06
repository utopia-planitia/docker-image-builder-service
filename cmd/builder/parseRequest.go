package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

func parseTagsAndBranches(r *http.Request) ([]*tag, []string, error) {

	tags := []*tag{}
	for _, t := range r.URL.Query()["t"] {
		tags = append(tags, newTag(t))
	}

	cfjson := r.URL.Query()["cachefrom"]
	if len(cfjson) == 0 {
		return tags, nil, nil
	}
	if len(cfjson) > 1 {
		return nil, nil, errors.New("cachefrom parameter is set multiple times")
	}
	cf, err := decodeCachefromJSON(cfjson[0])
	if err != nil {
		return nil, nil, fmt.Errorf("failed to decode cache from json: %s", err)
	}
	branches := filterBranches(cf)

	return tags, branches, nil
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

func filterBranches(cf []string) []string {
	var branches []string
	for _, e := range cf {
		if strings.HasPrefix(e, "branch=") {
			branches = append(branches, e[len("branch="):])
		}
	}
	return branches
}

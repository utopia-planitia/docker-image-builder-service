package main

import (
	"fmt"
	"log"
	"os/exec"
	"strings"
)

func load(tags []*tag, branches []string) {
	log.Printf("loading tags: %s / branches %s", tags, branches)
	if len(tags) != 0 {
		for _, t := range tags {
			loadCommand(t, cachedLatestFilename(t))
		}
	}

	for _, b := range branches {
		if b == "master" {
			continue
		}
		for _, t := range tags {
			loadCommand(t, cachedBranchFilename(t, b))
		}
	}
}

func save(tags []*tag, branch string) {
	log.Printf("saving tags: %s / currentBranch %s", tags, branch)
	for _, t := range tags {
		if t.version == "latest" {
			saveCommand(t, cachedLatestFilename(t))
		}
	}
	if branch == "master" {
		for _, t := range tags {
			saveCommand(t, cachedLatestFilename(t))
		}
	}
	for _, t := range tags {
		saveCommand(t, cachedBranchFilename(t, branch))
	}
}

func loadCommand(t fmt.Stringer, f filename) {
	log.Printf("loading cached file %s", f)
	/* #nosec */
	output, err := exec.Command("load", t.String(), string(f)).CombinedOutput()
	if err != nil {
		log.Printf("loading cached file %s failed: %v: %v", f, err, string(output))
	}
}

func saveCommand(t fmt.Stringer, f filename) {
	log.Printf("saving image %s to file %s", t, f)
	/* #nosec */
	output, err := exec.Command("save", t.String(), string(f)).CombinedOutput()
	if err != nil {
		log.Printf("saving image %s to file %s failed: %v: %v", t, f, err, string(output))
	}
}

func cachedLatestFilename(t *tag) filename {
	return filename(strings.Replace(t.image, "/", "~", -1) + ":latest")
}

func cachedBranchFilename(t *tag, branch string) filename {
	return filename(strings.Replace(t.image, "/", "~", -1) + ":" + branch)
}

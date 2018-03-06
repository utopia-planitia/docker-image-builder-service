package main

import (
	"fmt"
	"log"
	"os/exec"
	"strings"
)

func load(tags []*tag, branches []string) {
	if len(tags) != 0 {
		loadCommand(tags[0], cachedLatestFilename(tags[0]))
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

func save(tags []*tag, branches []string) {
	for _, t := range tags {
		if t.version != "latest" {
			continue
		}
		saveCommand(t, cachedLatestFilename(t))
	}
	if len(tags) == 0 {
		return
	}
	for _, b := range branches {
		if b == "master" {
			continue
		}
		saveCommand(tags[0], cachedBranchFilename(tags[0], b))
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

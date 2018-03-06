package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"net/http"
	"os/exec"
	"strconv"
	"strings"

	"code.cloudfoundry.org/bytefmt"
)

func uncachedBytes(w http.ResponseWriter, r *http.Request) {

	log.Printf("requested path: %s\n", r.URL)

	tags, branches, err := parseTagsAndBranches(r)
	if err != nil {
		log.Printf("failed to parse tags and branches from request: %s\n", err)
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	var bytes uint64

	if len(tags) != 0 {
		uncached, err := calculateUncachedBytes(tags[0], cachedLatestFilename(tags[0]))
		if err != nil {
			log.Printf("failed to calculate uncached bytes for tag %s: %s\n", tags[0], err)
			bytes = math.MaxInt64
		}
		if bytes != math.MaxInt64 {
			bytes += uncached
		}
	}

	for _, b := range branches {
		if b == "master" {
			continue
		}
		for _, t := range tags {
			uncached, err := calculateUncachedBytes(t, cachedBranchFilename(t, b))
			if err != nil {
				log.Printf("failed to calculate uncached bytes for branch %s / tag %s: %s\n", b, tags[0], err)
				bytes = math.MaxInt64
			}
			if bytes != math.MaxInt64 {
				bytes += uncached
			}
		}
	}

	_, err = w.Write([]byte(strconv.FormatUint(bytes, 10)))
	if err != nil {
		log.Printf("failed write result: %s\n", err)
	}
}

func uncachedBytesCacheFrom(cf []string, t *tag) (uint64, error) {
	var bytes uint64
	for _, e := range cf {
		f := cachedBranchFilename(t, e)
		b, err := calculateUncachedBytes(t, f)
		if err != nil {
			return 0, err
		}
		bytes += b
	}
	return bytes, nil
}

func calculateUncachedBytes(t *tag, f filename) (uint64, error) {
	/* #nosec */
	output, err := exec.Command("uncachedBytes", t.String(), string(f)).CombinedOutput()
	if err != nil {
		return 0, fmt.Errorf("crawling uncached file %s failed: %v: %v", f, err, string(output))
	}
	var bytes uint64
	scanner := bufio.NewScanner(strings.NewReader(string(output)))
	for scanner.Scan() {
		b, err := bytefmt.ToBytes(scanner.Text())
		if err != nil {
			return 0, fmt.Errorf("failed parse uncached bytes %v: %v", scanner.Text(), err)
		}
		bytes += b
	}
	if err := scanner.Err(); err != nil {
		return 0, fmt.Errorf("failed to scan uncached bytes: %v", err)
	}
	return bytes, nil
}

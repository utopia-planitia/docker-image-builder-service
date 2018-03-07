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

	tags, cacheFromBranches, _, err := parseTagsAndBranches(r)
	if err != nil {
		log.Printf("failed to parse tags and branches from request: %s\n", err)
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	var bytes uint64

	for _, t := range tags {
		f := cachedLatestFilename(t)
		bytes = addUncachedBytes(bytes, t, f)
	}

	for _, b := range cacheFromBranches {
		if b == masterBranch {
			continue
		}
		for _, t := range tags {
			f := cachedBranchFilename(t, b)
			bytes = addUncachedBytes(bytes, t, f)
		}
	}

	_, err = w.Write([]byte(strconv.FormatUint(bytes, 10)))
	if err != nil {
		log.Printf("failed write result: %s\n", err)
	}
}

func addUncachedBytes (bytes uint64, t fmt.Stringer, f filename) uint64 {
	uncached, err := calculateUncachedBytes(t, f)
	if err != nil {
		log.Printf("failed to calculate uncached bytes tag %s with filename %s: %s\n", t, f, err)
		bytes = math.MaxInt64
	}
	if bytes != math.MaxInt64 {
		bytes += uncached
	}
	return bytes
}

func calculateUncachedBytes(t fmt.Stringer, f filename) (uint64, error) {
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

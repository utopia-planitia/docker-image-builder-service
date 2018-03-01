package main

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"strconv"
	"strings"

	"code.cloudfoundry.org/bytefmt"
)

func uncachedBytes(w http.ResponseWriter, r *http.Request) {

	log.Printf("requested path: %s\n", r.URL)

	t, err := parseTag(r)
	if err != nil {
		log.Printf("failed to parse tag: %s\n", err)
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}
	cf, err := parseCachefromBranches(r)
	if err != nil {
		log.Printf("failed to parse branches: %s\n", err)
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	var bytes uint64

	f := cachedLatestFilename(t)
	b, err := calculateUncachedBytes(t, f)
	if err != nil {
		log.Printf("failed to add up uncached bytes: %s\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	bytes += b

	bytesCacheFrom, err := uncachedBytesCacheFrom(cf, t)
	if err != nil {
		log.Printf("failed to add up uncached bytes: %s\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	bytes += bytesCacheFrom

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

func calculateUncachedBytes(t fmt.Stringer, f filename) (uint64, error) {
	output, err := exec.Command("uncachedBytes", t.String(), string(f)).CombinedOutput()
	if err != nil {
		log.Printf("crawling uncached file %s failed: %v: %v", f, err, string(output))
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

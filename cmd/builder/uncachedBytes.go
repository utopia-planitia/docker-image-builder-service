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
	"github.com/damoon/docker-image-builder-service/dibs"
)

func uncachedBytes(w http.ResponseWriter, r *http.Request) {

	log.Printf("requested path: %s\n", r.URL)

	t, err := dibs.ParseTag(r)
	if err != nil {
		log.Printf("parameter t missing: %s\n", err)
		w.Write([]byte("parameter t missing"))
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}
	cf, err := parseCachefrom(r)
	if err != nil {
		log.Printf("parameter cachefrom missing: %s\n", err)
		w.Write([]byte("parameter cachefrom missing"))
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	var bytes uint64

	f := cachedLatestFilename(t)
	b, err := calculateUncachedBytes(t, f)
	if err != nil {
		log.Printf("failed to add up uncached bytes: %s\n", err)
		w.Write([]byte("failed to add up uncached bytes"))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	bytes += b

	for _, e := range cf {
		f = cachedBranchFilename(t, e)
		b, err := calculateUncachedBytes(t, f)
		if err != nil {
			log.Printf("failed to add up uncached bytes: %s\n", err)
			w.Write([]byte("failed to add up uncached bytes"))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		bytes += b
	}

	w.Write([]byte(strconv.FormatUint(bytes, 10)))
}

func calculateUncachedBytes(t *dibs.Tag, f filename) (uint64, error) {
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

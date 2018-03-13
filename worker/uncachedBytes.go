package main

// https://github.com/genuinetools/reg

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/docker/docker/api/types"
	"github.com/genuinetools/reg/registry"
	"github.com/genuinetools/reg/repoutils"
)

func uncachedBytes(w http.ResponseWriter, r *http.Request) {

	log.Printf("requested path: %s\n", r.URL)

	tags, branches, currentBranch, err := parseTagsAndBranches(r)
	if err != nil {
		log.Printf("failed to parse tags and branches from request: %s\n", err)
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	loadableTags := cacheSources(tags, branches, currentBranch)
	remoteLayers, err := remoteLayers(loadableTags)
	if err != nil {
		log.Printf("failed list remote layers: %s\n", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
	localLayers, err := localLayers()
	if err != nil {
		log.Printf("failed list local layers: %s\n", err)
		w.WriteHeader(http.StatusInternalServerError)
	}

	log.Printf("localLayers: %v\n", localLayers)
	log.Printf("remoteLayers: %v\n", remoteLayers)

	for _, layer := range localLayers {
		delete(remoteLayers, layer)
	}

	var bytes int64
	for _, size := range remoteLayers {
		bytes += size
	}

	_, err = w.Write([]byte(strconv.FormatInt(bytes, 10)))
	if err != nil {
		log.Printf("failed write result: %s\n", err)
	}
}

func remoteLayers(tags []*tag) (map[string]int64, error) {

	r, err := registry.NewInsecure(types.AuthConfig{ServerAddress: "http://cache:5000/"}, false)
	if err != nil {
		return nil, fmt.Errorf("failed setup registry client: %s", err)
	}

	layers := map[string]int64{}

	for _, tag := range tags {

		repo, ref, err := repoutils.GetRepoAndRef(tag.String())
		if err != nil {
			log.Printf("failed parse tag %s: %s\n", tag.String(), err)
			continue
		}

		manifest, err := r.ManifestV2(repo, ref)
		if err != nil {
			log.Printf("failed load manifest %s:%s: %s\n", repo, ref, err)
			continue
		}

		for _, layer := range manifest.Layers {
			log.Printf("remote layer %s has size %d\n", string(layer.Digest), layer.Size)
			layers[string(layer.Digest)] = layer.Size
		}
		log.Printf("remote layer %s has size %d\n", string(manifest.Config.Digest), manifest.Config.Size)
		layers[string(manifest.Config.Digest)] = manifest.Config.Size
	}
	return layers, nil
}

func localLayers() ([]string, error) {
	layers := []string{}
	files, err := ioutil.ReadDir("/var/lib/docker/image/overlay2/distribution/diffid-by-digest/sha256/")
	if err != nil {
		return nil, fmt.Errorf("failed list layers: %s", err)
	}
	for _, f := range files {
		layers = append(layers, "sha256:"+f.Name())
	}
	files, err = ioutil.ReadDir("/var/lib/docker/image/overlay2/imagedb/content/sha256/")
	if err != nil {
		return nil, fmt.Errorf("failed list layers: %s", err)
	}
	for _, f := range files {
		layers = append(layers, "sha256:"+f.Name())
	}
	return layers, nil
}

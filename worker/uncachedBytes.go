package main

// https://github.com/genuinetools/reg

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
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

	r, err := registry.NewInsecure(types.AuthConfig{}, false)
	if err != nil {
		return nil, fmt.Errorf("failed setup registry client: %s", err)
	}

	layers := map[string]int64{}

	for _, tag := range tags {

		repo, ref, err := repoutils.GetRepoAndRef(tag.String())
		if err != nil {
			log.Printf("failed load manifest %s:%s: %s\n", repo, ref, err)
			continue
		}

		manifest, err := r.ManifestV2(repo, ref)
		if err != nil {
			log.Printf("failed load manifest %s:%s: %s\n", repo, ref, err)
			continue
		}

		for _, layer := range manifest.Layers {
			layers[string(layer.Digest)] = layer.Size
		}
	}
	return layers, nil
}

func localLayers() ([]string, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return nil, fmt.Errorf("failed to create docker client: %s", err)
	}
	images, err := cli.ImageList(context.Background(), types.ImageListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed list images: %s", err)
	}
	layers := []string{}
	for _, image := range images {
		history, err := cli.ImageHistory(context.Background(), image.ID)
		if err != nil {
			log.Printf("failed get histroy of image %s: %s", image.ID, err)
		}
		for _, layer := range history {
			layers = append(layers, layer.ID)
		}
	}
	return layers, nil
}

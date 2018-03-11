package main

// https://github.com/moby/moby/tree/master/client#go-client-for-the-docker-engine-api

import (
	"context"
	"log"
	"strings"

	"bytes"

	"github.com/docker/distribution/reference"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

const masterBranch = "master"

func load(tags []*tag, branches []string, currentBranch string) {

	for _, t := range cacheSources(tags, branches, currentBranch) {
		loadCommand(t)
	}
}

func cacheSources(tags []*tag, branches []string, currentBranch string) []*tag {
	sources := []*tag{}

	log.Printf("loading tags: %s / branches %s", tags, branches)
	for _, t := range tags {
		log.Printf("loading :latest to tag %s", t)
		sources = append(sources, cachedLatestFilename(t))

		if currentBranch != "" {
			log.Printf("loading branch of tag %s / currentBranch %s", t, currentBranch)
			sources = append(sources, cachedBranchFilename(t, currentBranch))
		}
	}

	for _, b := range branches {
		if b == masterBranch {
			continue
		}
		for _, t := range tags {
			log.Printf("loading branch of tag %s / branch %s", t, currentBranch)
			sources = append(sources, cachedBranchFilename(t, b))
		}
	}

	return sources
}

func save(tags []*tag, branch string) {
	log.Printf("saving tags: %s / currentBranch %s", tags, branch)
	for _, t := range tags {
		if t.version == "latest" {
			log.Printf("saving :latest to tag %s", t)
			saveCommand(t, cachedLatestFilename(t))
		}
	}
	if branch == masterBranch {
		for _, t := range tags {
			log.Printf("saving masterBranch to tag %s", t)
			saveCommand(t, cachedLatestFilename(t))
		}
	}
	if branch != "" {
		for _, t := range tags {
			log.Printf("saving currentBranch to tag %s", t)
			saveCommand(t, cachedBranchFilename(t, branch))
		}
	}
}

func cachedLatestFilename(t *tag) *tag {
	return cachedBranchFilename(t, "latest")
}

func cachedBranchFilename(t *tag, branch string) *tag {
	return &tag{
		image:   "cache:5000/" + strings.Replace(t.image, ":", "~", -1),
		version: branch,
	}
}

func loadCommand(remote reference.Reference) {

	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		log.Printf("failed to create docker client: %s", err)
		return
	}

	log.Printf("pull image %s", remote.String())
	response, err := cli.ImagePull(context.Background(), remote.String(), types.ImagePullOptions{RegistryAuth: "Og=="})
	if err != nil {
		log.Printf("failed to pull image: %s", err)
		return
	}
	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(response)
	if err != nil {
		log.Printf("failed read from pull response: %s", err)
	}
	log.Printf("pull image response: %s", buf.String())

	err = response.Close()
	if err != nil {
		log.Printf("failed to close response of image pull: %s", err)
	}
}

func saveCommand(local, remote reference.Reference) {

	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		log.Printf("failed to create docker client: %s", err)
		return
	}

	log.Printf("tag image %s to %s", local.String(), remote.String())
	err = cli.ImageTag(context.Background(), local.String(), remote.String())
	if err != nil {
		log.Printf("failed to tag image %s to %s: %s", local, remote, err)
	}

	log.Printf("push image %s", remote.String())
	response, err := cli.ImagePush(context.Background(), remote.String(), types.ImagePushOptions{RegistryAuth: "Og=="})
	if err != nil {
		log.Printf("failed to push image: %s", err)
		return
	}
	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(response)
	if err != nil {
		log.Printf("failed read from push response: %s", err)
	}
	log.Printf("push image response: %s", buf.String())

	err = response.Close()
	if err != nil {
		log.Printf("failed to close response of image push: %s", err)
	}
}

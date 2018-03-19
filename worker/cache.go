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

func load(tag *tag, currentBranch string) {

	for _, t := range cacheSources(tag, currentBranch) {
		log.Printf("loading image %s as cache", t)
		loadCommand(t)
	}
}

func cacheSources(t *tag, currentBranch string) []*tag {
	// order of images is important
	// https://github.com/moby/moby/issues/26065#issuecomment-249046559
	// https://github.com/moby/moby/pull/26839#issuecomment-277383550

	sources := []*tag{}
	if currentBranch != "" {
		sources = append(sources, cachedBranchFilename(t, currentBranch))
	}
	sources = append(sources, cachedLatestFilename(t))
	return sources
}

func save(t *tag, branch string) {
	log.Printf("saving tags: %s / currentBranch %s", t, branch)
	if t.version == "latest" {
		log.Printf("saving :latest to tag %s", t)
		saveCommand(t, cachedLatestFilename(t))
	}
	if branch == masterBranch {
		log.Printf("saving masterBranch to tag %s", t)
		saveCommand(t, cachedLatestFilename(t))
	}
	if branch != "" && branch != masterBranch {
		log.Printf("saving currentBranch to tag %s", t)
		saveCommand(t, cachedBranchFilename(t, branch))
	}
}

func cachedLatestFilename(t *tag) *tag {
	return cachedBranchFilename(t, "latest")
}

func cachedBranchFilename(t *tag, branch string) *tag {
	return &tag{
		image:   strings.Replace(t.image, ":", "~", -1),
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
	response, err := cli.ImagePull(context.Background(), "cache:5000/"+remote.String(), types.ImagePullOptions{RegistryAuth: "Og=="})
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

	log.Printf("tag image %s to %s", local.String(), "cache:5000/"+remote.String())
	err = cli.ImageTag(context.Background(), local.String(), "cache:5000/"+remote.String())
	if err != nil {
		log.Printf("failed to tag image %s to %s: %s", local, remote, err)
	}

	log.Printf("push image %s", remote.String())
	response, err := cli.ImagePush(context.Background(), "cache:5000/"+remote.String(), types.ImagePushOptions{RegistryAuth: "Og=="})
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

package main

// https://github.com/moby/moby/tree/master/client#go-client-for-the-docker-engine-api

import (
	"bytes"
	"context"
	"io"
	"log"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

const masterBranch = "master"

func load(tag *tag, currentBranch string) {

	for _, t := range cacheSources(tag, currentBranch) {
		log.Printf("loading image %s as cache", t)
		loadCommand(t.String())
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

func loadCommand(remote string) {

	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		log.Printf("failed to create docker client: %s", err)
		return
	}

	log.Printf("pull image %s", remote)
	response, err := pull(cli, remote)
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

func pull(cli *client.Client, remote string) (io.ReadCloser, error) {
	ctx := context.Background()
	image := "cache:5000/" + remote
	options := types.ImagePullOptions{RegistryAuth: "Og=="}
	return cli.ImagePull(ctx, image, options)
}

func push(cli *client.Client, remote string) (io.ReadCloser, error) {
	ctx := context.Background()
	image := "cache:5000/" + remote
	options := types.ImagePushOptions{RegistryAuth: "Og=="}
	return cli.ImagePush(ctx, image, options)
}

func save(t *tag, branch string) {
	if t.version == "latest" || branch == masterBranch {
		log.Printf("saving latest cache image for tag %s", t)
		saveCommand(t.String(), cachedLatestFilename(t).String())
	}
	if branch != "" && branch != masterBranch {
		log.Printf("saving currentBranch to tag %s", t)
		saveCommand(t.String(), cachedBranchFilename(t, branch).String())
	}
}

func cachedLatestFilename(t *tag) *tag {
	return cachedBranchFilename(t, "latest")
}

func cachedBranchFilename(t *tag, branch string) *tag {
	return &tag{
		image:   strings.Replace(t.image, ":", "-", -1),
		version: branch,
	}
}

func saveCommand(local, remote string) {

	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		log.Printf("failed to create docker client: %s", err)
		return
	}

	log.Printf("tag image %s to %s", local, "cache:5000/"+remote)
	err = cli.ImageTag(context.Background(), local, "cache:5000/"+remote)
	if err != nil {
		log.Printf("failed to tag image %s to %s: %s", local, remote, err)
	}

	log.Printf("push image %s", remote)
	response, err := push(cli, remote)
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

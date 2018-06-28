package main

// https://github.com/moby/moby/tree/master/client#go-client-for-the-docker-engine-api

import (
	"bytes"
	"context"
	"log"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

func cacheSources(t *tag, currentBranch, headBranch string) []*tag {
	// order of images is important
	// https://github.com/moby/moby/issues/26065#issuecomment-249046559
	// https://github.com/moby/moby/pull/26839#issuecomment-277383550

	sources := []*tag{}
	sources = append(sources, cachedBranchImageTag(t, currentBranch))
	if currentBranch != headBranch {
		sources = append(sources, cachedBranchImageTag(t, headBranch))
	}
	return sources
}

func cachedBranchImageTag(t *tag, branch string) *tag {
	return &tag{
		image:   strings.Replace(t.image, ":", "-", -1),
		version: branch,
	}
}

func restoreCache(tag *tag, currentBranch, headBranch string) []*tag {
	for _, t := range cacheSources(tag, currentBranch, headBranch) {
		log.Printf("loading image %s as cache", t)
		pullFromCache(t.String())
	}
	return cacheSources(tag, currentBranch, headBranch)
}

func pullFromCache(remote string) {

	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		log.Printf("failed to create docker client: %s", err)
		return
	}

	log.Printf("pull image %s", remote)
	ctx := context.Background()
	image := "cache:5000/" + remote
	options := types.ImagePullOptions{RegistryAuth: "Og=="}
	response, err := cli.ImagePull(ctx, image, options)
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

func archiveCache(t *tag, currentBranch string) {
	archvie := cachedBranchImageTag(t, currentBranch).String()
	log.Printf("archive image %s to %s", t, archvie)
	pushToCache(t.String(), archvie)
}

func pushToCache(local, remote string) {

	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		log.Printf("failed to create docker client: %s", err)
		return
	}

	log.Printf("tag image %s to %s", local, "cache:5000/"+remote)
	err = cli.ImageTag(context.Background(), local, "cache:5000/"+remote)
	if err != nil {
		log.Printf("failed to tag image %s to %s: %s", local, remote, err)
		return
	}

	log.Printf("push image %s", remote)
	ctx := context.Background()
	image := "cache:5000/" + remote
	options := types.ImagePushOptions{RegistryAuth: "Og=="}
	response, err := cli.ImagePush(ctx, image, options)
	if err != nil {
		log.Printf("failed to push image: %s", err)
		return
	}
	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(response)
	if err != nil {
		log.Printf("failed read from push response: %s", err)
		return
	}
	log.Printf("push image response: %s", buf.String())

	err = response.Close()
	if err != nil {
		log.Printf("failed to close response of image push: %s", err)
	}
}

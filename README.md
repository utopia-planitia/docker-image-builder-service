# Docker image builder service

The scheduler of docker image builder service distributes docker builds to a set of workers.  
The workers contain a builder (local manager) and a docker in docker setup. The builder saves and loads the build history to minio.  
The scheduler forwards to the same worker based on the client ip for 10s.  
The worker with the least amount of cache to catch up will be preferrd.  

Build [![CircleCI](https://circleci.com/gh/utopia-planitia/docker-image-builder-service.svg?style=svg)](https://circleci.com/gh/utopia-planitia/docker-image-builder-service)

[![Go Report Card](https://goreportcard.com/badge/github.com/utopia-planitia/docker-image-builder-service)](https://goreportcard.com/report/github.com/utopia-planitia/docker-image-builder-service)

## Development

local development happens via minikube
```
make minikube-start
make minikube-deploy
make minikube-test
make minikube-stop
```

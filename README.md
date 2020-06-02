# docker image builder

The scheduler of docker image builder service distributes docker builds to a set of workers.  
The workers contain a builder (local manager) and a docker in docker setup. The builder saves and loads the build history to minio.  
The scheduler forwards to the same worker based on the client ip for 10s.  
The worker with the least amount of cache to catch up will be preferrd.  

[![CircleCI](https://circleci.com/gh/utopia-planitia/docker-image-builder.svg?style=shield)](https://circleci.com/gh/utopia-planitia/docker-image-builder)

[![Go Report Card](https://goreportcard.com/badge/github.com/utopia-planitia/docker-image-builder)](https://goreportcard.com/report/github.com/utopia-planitia/docker-image-builder)

## Development

Development happens in kubernetes via tilt.
```
tilt up
```

Some tests and linter run locally.
```
make test
make lint
```

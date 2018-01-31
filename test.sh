#!/bin/bash

set -e

docker save -o alpine37.tar alpine:3.7

export DOCKER_HOST=tcp://127.0.0.1:2376
docker version

export DOCKER_HOST=tcp://127.0.0.1:2375
docker version

#docker pull alpine:3.7
docker load -i alpine37.tar

date

docker build test --no-cache=true -t a1 &

date

time wait

date

#docker build test --no-cache=true -t b1 &
#docker build test --no-cache=true -t b2 &
#docker build test --no-cache=true -t b3 &
#docker build test --no-cache=true -t b4 &
#docker build test --no-cache=true -t b5 &

#date

#time wait

unset DOCKER_HOST

#!/bin/bash

set -euo pipefail
#set -x

START=$(date +%s)

for _ in {1..3}
do
    DATE=$(date +%s%N)
    docker run --rm \
        --dns 10.96.0.10 \
        --dns-search container-image-builder.svc.cluster.local \
        -e DOCKER_HOST=tcp://docker:2375 \
        -v "$(pwd):/project" -w /project \
        utopiaplanitia/docker-image-builder-devtools:latest \
        docker build --build-arg version="$DATE" tests/example-build \
        &
done

wait
END=$(date +%s)

DELAY=$((END - START))
echo "build took $DELAY seconds"
if [ "$DELAY" -lt 9 ]; then
  echo "build was to fast"
  exit 1
fi
if [ "$DELAY" -gt 30 ]; then
  echo "build was to slow"
  exit 2
fi

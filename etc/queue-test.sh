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
        -e CACHE_ENDPOINT_SERVER=minio \
        -e CACHE_ENDPOINT_PORT=9000 \
        -e CACHE_BUCKET=image-layers \
        -e CACHE_ACCESS_KEY=8Q9U4RBHKKB6HU70SRZ1 \
        -e CACHE_SECRET_KEY=oxxT2iqBlW6lgaDVe8ll6mP8z/OSVIUnn9cB4+Q0 \
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

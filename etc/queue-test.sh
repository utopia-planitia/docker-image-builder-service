#!/bin/bash

set -euo pipefail
#set -x

echo wait for workers to cool down
sleep 10
echo starting 5 builds

START=$(date +%s)

for i in {1..5}
do
    DATE=$(date +%s%N)
    docker run --rm \
        --dns 10.96.0.10 \
        --dns-search container-image-builder.svc.cluster.local \
        -e DOCKER_HOST=tcp://docker:2375 \
        -v "$(pwd):/project" -w /project \
        utopiaplanitia/docker-image-builder-devtools:latest \
        docker build --build-arg version="$DATE" tests/example-build \
        > $i.log \
        &
done

wait
END=$(date +%s)

DELAY=$((END - START))
echo "build took $DELAY seconds"
if [ "$DELAY" -lt 20 ]; then
  echo "build was to fast"
  for i in {1..5}
  do
      echo "cat $i.log"
      cat $i.log
  done
  exit 1
fi
if [ "$DELAY" -gt 40 ]; then
  echo "build was to slow"
  exit 2
fi

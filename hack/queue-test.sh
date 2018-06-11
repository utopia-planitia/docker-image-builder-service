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
        -v "$(pwd)/hack/queue-test":/root/.docker \
        utopiaplanitia/docker-image-builder-devtools:latest \
        docker build --build-arg version="$DATE" -t "queue-test-$DATE" tests/example-build \
        > $i.log \
        &
done

wait
END=$(date +%s)

DELAY=$((END - START))
echo "build took $DELAY seconds"

if [ "$DELAY" -gt 10 ]; then
  echo "build was to slow"
  for i in {1..5}
  do
      echo "cat $i.log"
      cat $i.log
  done
  exit 2
fi

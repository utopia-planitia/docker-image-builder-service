#!/bin/bash

# $1 = image
# $2 = file

if ! mc ls "cache/$CACHE_BUCKET" > /dev/null; then
  mc mb "cache/$CACHE_BUCKET"
fi

S3_HASH=$(mc cat "cache/$CACHE_BUCKET/$2.meta" | tail -n 1)
LOCAL_HASH=$(docker history -q "$1" | md5sum | head -c 32)
if [ "$S3_HASH" == "$LOCAL_HASH" ]; then
  echo "hashes match, skip build cache save"
  exit 11
fi

set -e
LAYERS=$(docker history -q "$1" | grep -v missing)
# shellcheck disable=SC2086
docker save "$1" $LAYERS | mc pipe "cache/$CACHE_BUCKET/$2"
DATE=$(date +%s)
HASH=$(docker history -q "$1" | md5sum | head -c 32)
echo -e "$DATE\n$HASH" | mc pipe "cache/$CACHE_BUCKET/$2.meta"

set +e
echo -n waiting for meta data
until mc stat "cache/$CACHE_BUCKET/$2.meta"; do
    sleep 1
    echo -n "."
done
echo
echo -n waiting for data file
until mc stat "cache/$CACHE_BUCKET/$2" | grep -v meta | grep "Name      :"; do
    sleep 1
    echo -n "."
done
echo

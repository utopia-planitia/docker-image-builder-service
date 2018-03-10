#!/bin/bash

# $1 = image
# $2 = file

set -e

S3_HASH=$(mc cat "cache/$CACHE_BUCKET/$2.meta" | tail -n 1)
if [ "${S3_HASH}" == "" ]; then
  echo "meta file in s3 is missing"
  exit 12
fi
LOCAL_HASH=$(docker history -q "$1" | md5sum | head -c 32)
if [ "$S3_HASH" == "$LOCAL_HASH" ]; then
  echo "hashes match, skip build cache load"
  exit 11
fi

mc cat "cache/$CACHE_BUCKET/$2" | docker load

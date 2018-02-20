#!/bin/sh

# $1 = image
# $2 = file

S3_HASH=$(mc cat cache/$CACHE_BUCKET/$2 | tail -n 1)
LOCAL_HASH=$(docker history -q $1 | md5sum | head -c 32)
if [ "${S3_HASH}" -eq "${LOCAL_HASH}" ]; then
  echo "hashes match, skip build cache load"
  exit 11
fi

mc cat cache/$CACHE_BUCKET/$2 | docker load

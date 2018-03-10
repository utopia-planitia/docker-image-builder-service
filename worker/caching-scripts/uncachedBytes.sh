#!/bin/bash

# $1 = image
# $2 = file

S3_HASH=$(mc cat "cache/$CACHE_BUCKET/$2.meta" | tail -n 1)
LOCAL_HASH=$(docker history -q "$1" | md5sum | head -c 32)
if [ "$S3_HASH" == "$LOCAL_HASH" ]; then
  echo "0B"
  exit
fi

mc stat "cache/$CACHE_BUCKET/$2" | grep Size | awk '{print $3}'

#!/bin/sh

# $1 = image
# $2 = file

mc config host ls | grep $CACHE_SECRET_KEY > /dev/null
if [ $? != 0 ]; then
  mc config host add cache $CACHE_ENDPOINT $CACHE_ACCESS_KEY $CACHE_SECRET_KEY
fi

S3_HASH=$(mc cat cache/$CACHE_BUCKET/$2 | tail -n 1)
LOCAL_HASH=$(docker history -q $1 | md5sum)
if [ "${S3_HASH}" -eq "${LOCAL_HASH}" ]; then
  echo "hashes match, skip build cache load"
  exit 11
fi

mc cat cache/$CACHE_BUCKET/$2 | docker load

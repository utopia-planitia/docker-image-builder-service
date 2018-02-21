#!/bin/sh

# $1 = image
# $2 = file

mc ls cache/$CACHE_BUCKET > /dev/null
if [ $? != 0 ]; then
  mc mb cache/$CACHE_BUCKET
fi

set -e
DATE=$(date +%s)
HASH=$(docker history -q $1 | md5sum | head -c 32)
echo -e "$DATE\n$HASH" | mc pipe cache/$CACHE_BUCKET/$2.meta
docker save $1 $(docker history -q $1 | grep -v missing) | mc pipe cache/$CACHE_BUCKET/$2

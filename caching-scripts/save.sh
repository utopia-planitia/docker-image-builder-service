#!/bin/sh

# $1 = image
# $2 = file

mc config host ls | grep $CACHE_SECRET_KEY > /dev/null
if [ $? != 0 ]; then
  mc config host add cache $CACHE_ENDPOINT $CACHE_ACCESS_KEY $CACHE_SECRET_KEY
fi

mc ls cache/$CACHE_BUCKET > /dev/null
if [ $? != 0 ]; then
  mc mb cache/$CACHE_BUCKET
fi

docker save $1 $(docker history -q $1 | grep -v missing) | mc pipe cache/$CACHE_BUCKET/$2

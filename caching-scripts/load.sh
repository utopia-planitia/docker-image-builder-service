#!/bin/sh

# $1 = file

mc config host ls | grep $CACHE_SECRET_KEY > /dev/null
if [ $? != 0 ]; then
  mc config host add cache $CACHE_ENDPOINT $CACHE_ACCESS_KEY $CACHE_SECRET_KEY
fi

mc cat cache/$CACHE_BUCKET/$1 | docker load

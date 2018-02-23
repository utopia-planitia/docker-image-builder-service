#!/bin/bash

mc config host ls | grep $CACHE_SECRET_KEY > /dev/null
if [ $? != 0 ]; then
  mc config host add cache $CACHE_ENDPOINT $CACHE_ACCESS_KEY $CACHE_SECRET_KEY
fi

exec "$@"
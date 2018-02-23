#!/bin/bash

TIMEOUT="timeout"
$TIMEOUT 1 true
if [ $? != 0 ]; then
  TIMEOUT="timeout -t"
fi

$TIMEOUT 15 bash -c "until echo > /dev/tcp/$CACHE_ENDPOINT_SERVER/$CACHE_ENDPOINT_PORT; do sleep 0.5; done"
if [ $? != 0 ]; then
  echo s3 endpoint is not available
  exit 1
fi

mc config host ls | grep $CACHE_SECRET_KEY > /dev/null
if [ $? != 0 ]; then
  mc config host add cache http://$CACHE_ENDPOINT_SERVER:$CACHE_ENDPOINT_PORT/ $CACHE_ACCESS_KEY $CACHE_SECRET_KEY
fi

exec "$@"
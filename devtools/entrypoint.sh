#!/bin/bash

if ! mc config host ls | grep "$CACHE_SECRET_KEY" > /dev/null; then
  until mc config host add cache "http://$CACHE_ENDPOINT_SERVER:$CACHE_ENDPOINT_PORT/" "$CACHE_ACCESS_KEY" "$CACHE_SECRET_KEY"; do
    sleep 0.5;
  done
fi

exec "$@"

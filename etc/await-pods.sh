#!/bin/bash

sleep 1

echo -n waiting for minio
for _ in {1..150}; do # timeout for 5 minutes
  if kubectl -n container-image-builder -l app=minio get po | grep Running | grep 1/1 > /dev/null 2>&1; then
      echo " done"
      break
  fi
  echo -n .
  sleep 2
done

echo -n waiting for mirror
for _ in {1..150}; do # timeout for 5 minutes
  if kubectl -n container-image-builder -l app=mirror get po | grep Running | grep 1/1 > /dev/null 2>&1; then
      echo " done"
      break
  fi
  echo -n .
  sleep 2
done

echo -n waiting for dispatcher
for _ in {1..150}; do # timeout for 5 minutes
  if kubectl -n container-image-builder -l app=dispatcher get po | grep Running | grep 1/1 > /dev/null 2>&1; then
      echo " done"
      break
  fi
  echo -n .
  sleep 2
done

for i in {0..1}; do
  echo -n "waiting for builder-$i"
  for _ in {1..150}; do # timeout for 5 minutes
    if kubectl -n container-image-builder get po "builder-$i" | grep Running | grep 2/2 > /dev/null 2>&1; then
        echo " done"
        break
    fi
    echo -n .
    sleep 2
  done
done

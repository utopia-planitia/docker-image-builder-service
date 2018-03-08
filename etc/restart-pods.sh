#!/bin/bash

kubectl -n container-image-builder delete po --all

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

echo -n waiting for dispatcher
for _ in {1..150}; do # timeout for 5 minutes
  if kubectl -n container-image-builder -l app=dispatcher get po | grep Running | grep 1/1 > /dev/null 2>&1; then
      echo " done"
      break
  fi
  echo -n .
  sleep 2
done

echo -n waiting for builder-0
for _ in {1..150}; do # timeout for 5 minutes
  if kubectl -n container-image-builder get po builder-0 | grep Running | grep 2/2 > /dev/null 2>&1; then
      echo " done"
      break
  fi
  echo -n .
  sleep 2
done

echo -n waiting for builder-1
for _ in {1..150}; do # timeout for 5 minutes
  if kubectl -n container-image-builder get po builder-1 | grep Running | grep 2/2 > /dev/null 2>&1; then
      echo " done"
      break
  fi
  echo -n .
  sleep 2
done

echo -n waiting for dns to point to minio
for _ in {1..150}; do # timeout for 5 minutes
  if dig minio.container-image-builder.svc.cluster.local @10.96.0.10 | grep "ANSWER SECTION" > /dev/null 2>&1; then
      echo " done"
      break
  fi
  echo -n .
  sleep 2
done

echo -n waiting for dns to point to docker
for _ in {1..150}; do # timeout for 5 minutes
  if dig docker.container-image-builder.svc.cluster.local @10.96.0.10 | grep "ANSWER SECTION" > /dev/null 2>&1; then
      echo " done"
      break
  fi
  echo -n .
  sleep 2
done

echo -n waiting for dns to point to builder-0
for _ in {1..150}; do # timeout for 5 minutes
  if dig builder-0.builder.container-image-builder.svc.cluster.local @10.96.0.10 | grep "ANSWER SECTION" > /dev/null 2>&1; then
      echo " done"
      break
  fi
  echo -n .
  sleep 2
done

echo -n waiting for dns to point to builder-1
for _ in {1..150}; do # timeout for 5 minutes
  if dig builder-1.builder.container-image-builder.svc.cluster.local @10.96.0.10 | grep "ANSWER SECTION" > /dev/null 2>&1; then
      echo " done"
      break
  fi
  echo -n .
  sleep 2
done

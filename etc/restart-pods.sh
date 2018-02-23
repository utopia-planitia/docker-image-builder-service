#!/bin/bash

kubectl -n container-image-builder -l app=builder delete po
kubectl -n container-image-builder -l app=dispatcher delete po

sleep 1

echo -n waiting for dispatcher
until kubectl -n container-image-builder -l app=dispatcher get po | grep Running | grep 1/1 > /dev/null; do
    sleep 1
    echo -n "."
done
echo

echo -n waiting for builder-0
until kubectl -n container-image-builder get po builder-0 | grep Running | grep 2/2 > /dev/null; do
    sleep 1
    echo -n "."
done
echo

echo -n waiting for builder-1
until kubectl -n container-image-builder get po builder-1 | grep Running | grep 2/2 > /dev/null; do
    sleep 1
    echo -n "."
done
echo

#!/bin/bash

kubectl -n container-image-builder delete po --all

source ./etc/await-pods.sh

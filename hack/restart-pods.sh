#!/bin/bash

kubectl -n container-image-builder delete po --all

source ./hack/await-pods.sh

#!/bin/bash

kubectl -n container-image-builder delete po --all

source await-pods.sh

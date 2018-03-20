#!/bin/bash

echo wait for docker version
export DOCKER_HOST="tcp://docker:2375"
until docker version > /dev/null; do sleep 1; done

echo wait for builder 1 version
export DOCKER_HOST="tcp://builder-0.builder:2375"
until docker version > /dev/null; do sleep 1; done

echo wait for builder 2 version
export DOCKER_HOST="tcp://builder-1.builder:2375"
until docker version > /dev/null; do sleep 1; done

echo wait for mirror port
until curl --fail http://mirror:5000/v2/ > /dev/null 2> /dev/null; do sleep 1; done

echo wait for cache port
until curl --fail http://cache:5000/v2/ > /dev/null 2> /dev/null; do sleep 1; done

#!/bin/bash

# Define directory structure
SBIN_DIR=`dirname $0`
SBIN_DIR="`cd "$SBIN_DIR"; pwd`"
DOCKER_DIR=$SBIN_DIR/../../docker

# Spin up containers
docker compose -f $DOCKER_DIR/docker-compose-integration-test.yml build
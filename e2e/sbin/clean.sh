#!/bin/bash

# Define directory structure
SBIN_DIR=`dirname $0`
SBIN_DIR="`cd "$SBIN_DIR"; pwd`"
PROJECT_DIR=$SBIN_DIR/../project
DOCKER_DIR=$SBIN_DIR/../../docker

# Clean up
cd $PROJECT_DIR
docker compose -f $DOCKER_DIR/docker-compose-integration-test.yml down
rm -rf $PROJECT_DATA_DIR/geth $PROJECT_DATA_DIR/keystore $PROJECT_DATA_DIR/geth.ipc

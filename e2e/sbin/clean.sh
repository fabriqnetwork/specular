#!/bin/bash

# Define directory structure
SBIN_DIR=`dirname $0`
SBIN_DIR="`cd "$SBIN_DIR"; pwd`"
PROJECT_DIR=$SBIN_DIR/../project
PROJECT_DATA_DIR=$PROJECT_DIR/specular-datadir
DOCKER_DIR=$SBIN_DIR/../../docker

# Clean up
cd $PROJECT_DIR
docker compose -f $DOCKER_DIR/docker-compose-integration-test.yml down
rm -rf $PROJECT_DATA_DIR/geth
rm -rf $PROJECT_DATA_DIR/keystore
rm -rf $PROJECT_DATA_DIR/geth.ipc

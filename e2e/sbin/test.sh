#!/bin/bash

# Define constants
HOST=localhost
L1_PORT=8545
L2_PORT=4011

# Define directory structure
SBIN_DIR=`dirname $0`
SBIN_DIR="`cd "$SBIN_DIR"; pwd`"
PROJECT_DIR=$SBIN_DIR/../project
PROJECT_DATA_DIR=$PROJECT_DIR/specular-datadir
CONTRACTS_DIR=$SBIN_DIR/../../contracts
DOCKER_DIR=$SBIN_DIR/../../docker

# Spin up containers
cd $PROJECT_DIR
docker compose -f $DOCKER_DIR/docker-compose-integration-test.yml up -d
sleep 30
$SBIN_DIR/wait-for-it.sh -t 240 $HOST:$L1_PORT
$SBIN_DIR/wait-for-it.sh -t 240 $HOST:$L2_PORT

# Run testing script
cd $CONTRACTS_DIR
npx ts-node scripts/testing.ts
RESULT=$?

# Clean up
$SBIN_DIR/clean.sh

# Exit with result
exit $RESULT
#!/bin/bash

# Configure variables
SBIN_DIR=`dirname $0`
SBIN_DIR="`cd "$SBIN_DIR"; pwd`"
set -o allexport
source $SBIN_DIR/configure.sh
set +o allexport

docker remove --force geth_container

# Clean up
rm -rf $PROJECT_DATA_DIR/geth
rm -rf $PROJECT_DATA_DIR/keystore
rm -rf $PROJECT_DATA_DIR/geth.ipc

# Remove deployments
rm -rf $CONTRACTS_DIR/deployments/localhost
rm -rf $CONTRACTS_DIR/deployments/specularLocalDev


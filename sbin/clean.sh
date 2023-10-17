#!/bin/bash
SBIN=`dirname $0`
SBIN="`cd "$SBIN"; pwd`"
. $SBIN/configure.sh

# Remove L1 docker container
docker remove --force geth_container

# Clean up data dir
rm -rf $DATA_DIR/geth
rm -rf $DATA_DIR/keystore
rm -rf $DATA_DIR/geth.ipc

# Remove deployments
rm -rf $CONTRACTS_DIR/deployments/localhost
rm -rf $CONTRACTS_DIR/deployments/specularLocalDev

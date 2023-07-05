#!/bin/bash

# Configure variables
SBIN_DIR=`dirname $0`
SBIN_DIR="`cd "$SBIN_DIR"; pwd`"
set -o allexport
source $SBIN_DIR/configure.sh
set +o allexport

# Clean up
rm -rf $PROJECT_DATA_DIR/geth
rm -rf $PROJECT_DATA_DIR/keystore
rm -rf $PROJECT_DATA_DIR/geth.ipc

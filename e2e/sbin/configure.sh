#!/bin/bash

# Define directory structure
SBIN_DIR=`dirname $0`
SBIN_DIR="`cd "$SBIN_DIR"; pwd`"
DATA_DIR=$SBIN_DIR/../data
PROJECT_DIR=$SBIN_DIR/../project
PROJECT_LOG_DIR=$PROJECT_DIR/logs
PROJECT_DATA_DIR=$PROJECT_DIR/specular-datadir
CONFIG_DIR=$SBIN_DIR/../../config
CONTRACTS_DIR=$SBIN_DIR/../../contracts
GETH_SPECULAR_DIR=$SBIN_DIR/../../services/sidecar
L2GETH_BIN=$GETH_SPECULAR_DIR/build/bin/geth

# Load environment variables
source $DATA_DIR/e2e.env

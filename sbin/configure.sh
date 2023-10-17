#!/bin/bash
# Define directory structure for other scripts.
SBIN=$(dirname "$(readlink -f "$0")")
SBIN="`cd "$SBIN"; pwd`"

ROOT_DIR=$SBIN/..

CONTRACTS_DIR=$ROOT_DIR/contracts
DATA_DIR=$ROOT_DIR/e2e/data
CONFIG_DIR=$ROOT_DIR/config
GETH_DIR=$ROOT_DIR/services/el_clients/go-ethereum
SIDECAR_DIR=$ROOT_DIR/services/sidecar

# Define binaries
SIDECAR_BIN=$SIDECAR_DIR/build/bin/sidecar
GETH_BIN=$GETH_DIR/build/bin/geth

# Load environment variables
. $DATA_DIR/e2e.env

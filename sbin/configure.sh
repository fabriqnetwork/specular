#!/bin/bash
# Define directory structure for other scripts.
# TODO: make this a dotenv
SBIN=$(dirname "$(readlink -f "$0")")
SBIN="`cd "$SBIN"; pwd`"
ROOT_DIR=$SBIN/..

CONTRACTS_DIR=$ROOT_DIR/contracts
CONFIG_DIR=$ROOT_DIR/config
GETH_DIR=$ROOT_DIR/services/el_clients/go-ethereum
MAGI_DIR=$ROOT_DIR/services/cl_clients/magi
SIDECAR_DIR=$ROOT_DIR/services/sidecar
# Define binaries
SP_GETH_BIN=$GETH_DIR/build/bin/geth
SP_MAGI_BIN=$MAGI_DIR/target/debug/magi
SIDECAR_BIN=$SIDECAR_DIR/build/bin/sidecar

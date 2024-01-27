#!/bin/bash

set -e

# Get the current script's directory
SCRIPT_DIR=$(dirname "$(readlink -f "$0")")
ROOT_DIR=$SCRIPT_DIR/..

# Load environment variables from .paths.env
PATHS_ENV=".paths.env"
if ! test -f "$PATHS_ENV"; then
  echo "Error: $PATHS_ENV not found"
  exit
fi
echo "Using dotenv: $PATHS_ENV"
. $PATHS_ENV

# Load environment variables from .genesis.env if it exists
GENESIS_ENV=".genesis.env"
if test -f "$GENESIS_ENV"; then
  . $GENESIS_ENV
fi

# Clean up existing files if they exist
cleanup_file() {
  if test -f "$1"; then
    echo "Removing $1"
    rm $1
  fi
}

cleanup_file "$GENESIS_PATH"
cleanup_file "$GENESIS_EXPORTED_HASH_PATH"
cleanup_file "$ROLLUP_CFG_PATH"
cleanup_file "$DEPLOYMENTS_ENV"

# Remove deployment files and .deployed file
echo "Removing deployment files in $CONTRACTS_DIR/deployments/$L1_NETWORK"
rm -rf $CONTRACTS_DIR/deployments/$L1_NETWORK
rm -f .deployed

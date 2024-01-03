#!/bin/bash
SBIN=$(dirname "$(readlink -f "$0")")
ROOT_DIR=$SBIN/..

PATHS_ENV=".paths.env"
if ! test -f "$PATHS_ENV"; then
  echo "Expected dotenv at $PATHS_ENV (does not exist)."
  exit
fi
echo "Using dotenv: $PATHS_ENV"
. $PATHS_ENV

GENESIS_ENV=".genesis.env"
if test -f "$GENESIS_ENV"; then
  . $GENESIS_ENV
fi

if test -f "$GENESIS_PATH"; then
  echo "Removing $GENESIS_PATH"
  rm $GENESIS_PATH
fi
if test -f "$GENESIS_EXPORTED_HASH_PATH"; then
  echo "Removing $GENESIS_EXPORTED_HASH_PATH"
  rm $GENESIS_EXPORTED_HASH_PATH
fi
if test -f "$ROLLUP_CFG_PATH"; then
  echo "Removing $ROLLUP_CFG_PATH"
  rm $ROLLUP_CFG_PATH
fi
DEPLOYMENTS_ENV=".deployments.env"
if test -f "$DEPLOYMENTS_ENV"; then
  echo "Removing $DEPLOYMENTS_ENV"
  rm $DEPLOYMENTS_ENV
fi

echo "Removing deployment files in $CONTRACTS_DIR/deployments/$L1_NETWORK"
rm -rf $CONTRACTS_DIR/deployments/$L1_NETWORK
rm -f .deployed

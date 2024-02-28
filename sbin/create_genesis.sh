#!/bin/bash
set -e

WORKSPACE_DIR=$HOME/.spc/workspaces/active_workspace

# Currently the local sbin paths are relative to the project root.
SBIN=$(dirname "$(readlink -f "$0")")
SBIN="$(
  cd "$SBIN"
  pwd
)"
. $SBIN/utils/utils.sh
ROOT_DIR=$SBIN/..

# Check that the all required dotenv files exists.
PATHS_ENV=$WORKSPACE_DIR/.paths.env
GENESIS_ENV=$WORKSPACE_DIR/.genesis.env
CONTRACTS_ENV=$WORKSPACE_DIR/.contracts.env
DEPLOYMENTS_ENV=$WORKSPACE_DIR/.deployments.env

reqdotenv "paths" $PATHS_ENV
reqdotenv "genesis" $GENESIS_ENV
reqdotenv "contracts" $CONTRACTS_ENV
reqdotenv "deployments" $DEPLOYMENTS_ENV

echo "Using $OPS_DIR as ops directory."
# Get relative paths for $OPS_DIR
GENESIS_CFG_PATH=$(relpath $GENESIS_CFG_PATH $OPS_DIR)
GENESIS_PATH=$(relpath $GENESIS_PATH $OPS_DIR)
GENESIS_EXPORTED_HASH_PATH=$(relpath $GENESIS_EXPORTED_HASH_PATH $OPS_DIR)
echo "Generating new genesis file at $GENESIS_PATH and exporting hash to $GENESIS_EXPORTED_HASH_PATH"
cd $OPS_DIR
guard_overwrite $GENESIS_PATH $AUTO_ACCEPT
# Create genesis.json file.
FLAGS=(
  "--genesis-config $GENESIS_CFG_PATH"
  "--out $GENESIS_PATH"
  "--l1-block $SEQUENCER_INBOX_DEPLOYED_BLOCK"
  "--l1-rpc-url $L1_ENDPOINT"
  "--export-hash $GENESIS_EXPORTED_HASH_PATH"
  "--l1-portal-address $L1PORTAL_ADDR"
  "--l1-standard-bridge-address $L1STANDARD_BRIDGE_ADDR"
  "--alloc $SEQUENCER_ADDRESS,$VALIDATOR_ADDRESS,$DEPLOYER_ADDRESS"
)

# hoop: I don't have the patience rn to determine why this isn't being sourced
if [[ -z ${OPS_GENESIS_BIN+x} ]]; then
  CMD="/usr/local/bin/genesis ${FLAGS[@]}"
else
  CMD="$OPS_GENESIS_BIN ${FLAGS[@]}"
fi

echo "Running $CMD"
eval $CMD

# Initialize a reference to the config files at
# "contracts/.genesis" (using relative paths as appropriate).
CONTRACTS_ENV=$CONTRACTS_DIR/.genesis.env
guard_overwrite $CONTRACTS_ENV $AUTO_ACCEPT
# Write file, using relative paths.
echo "Initializing contracts genesis dotenv $CONTRACTS_ENV"
GENESIS_PATH=$(relpath $GENESIS_PATH $CONTRACTS_DIR)
GENESIS_EXPORTED_HASH_PATH=$(relpath $GENESIS_EXPORTED_HASH_PATH $CONTRACTS_DIR)
BASE_ROLLUP_CFG_PATH=$(relpath $BASE_ROLLUP_CFG_PATH $CONTRACTS_DIR)
echo GENESIS_PATH=$GENESIS_PATH >>$CONTRACTS_ENV
echo GENESIS_EXPORTED_HASH_PATH=$GENESIS_EXPORTED_HASH_PATH >>$CONTRACTS_ENV

#!/bin/bash
set -e

# Set the SBIN variable to the directory of the script
SBIN=$(dirname "$(readlink -f "$0")")
SBIN="$(
  cd "$SBIN"
  pwd
)"
. $SBIN/utils/utils.sh
ROOT_DIR=$SBIN/..

# Check that all required dotenv files exist
require_dotenv "paths" ".paths.env"
require_dotenv "genesis" ".genesis.env"
require_dotenv "contracts" ".contracts.env"
require_dotenv "deployments" ".deployments.env"

# Output the ops directory being used
echo "Using $OPS_DIR as ops directory."

# Get relative paths for $OPS_DIR
GENESIS_CFG_PATH=$(relative_path $GENESIS_CFG_PATH $OPS_DIR)
GENESIS_PATH=$(relative_path $GENESIS_PATH $OPS_DIR)
GENESIS_EXPORTED_HASH_PATH=$(relative_path $GENESIS_EXPORTED_HASH_PATH $OPS_DIR)

# Generate new genesis file and export hash
echo "Generating new genesis file at $GENESIS_PATH and exporting hash to $GENESIS_EXPORTED_HASH_PATH"
cd $OPS_DIR
confirm_overwrite $GENESIS_PATH $AUTO_ACCEPT

# Create genesis.json file
FLAGS=(
  "--genesis-config $GENESIS_CFG_PATH"
  "--out $GENESIS_PATH"
  "--l1-rpc-url $L1_ENDPOINT"
  "--export-hash $GENESIS_EXPORTED_HASH_PATH"
  "--l1-portal-address $L1PORTAL_ADDR"
  "--l1-standard-bridge-address $L1STANDARD_BRIDGE_ADDR"
  "--alloc $SEQUENCER_ADDRESS,$VALIDATOR_ADDRESS,$DEPLOYER_ADDRESS"
)

# Choose the correct genesis binary and run it
if [[ -z ${OPS_GENESIS_BIN+x} ]]; then
  CMD="/usr/local/bin/genesis ${FLAGS[@]}"
else
  CMD="$OPS_GENESIS_BIN ${FLAGS[@]}"
fi

echo "Running $CMD"
eval $CMD

# Initialize a reference to the config files at "contracts/.genesis" using relative paths as appropriate
CONTRACTS_ENV=$CONTRACTS_DIR/.genesis.env
confirm_overwrite $CONTRACTS_ENV $AUTO_ACCEPT

# Write file using relative paths
echo "Initializing contracts genesis dotenv $CONTRACTS_ENV"
GENESIS_PATH=$(relative_path $GENESIS_PATH $CONTRACTS_DIR)
GENESIS_EXPORTED_HASH_PATH=$(relative_path $GENESIS_EXPORTED_HASH_PATH $CONTRACTS_DIR)
BASE_ROLLUP_CFG_PATH=$(relative_path $BASE_ROLLUP_CFG_PATH $CONTRACTS_DIR)
echo GENESIS_PATH=$GENESIS_PATH >>$CONTRACTS_ENV
echo GENESIS_EXPORTED_HASH_PATH=$GENESIS_EXPORTED_HASH_PATH >>$CONTRACTS_ENV

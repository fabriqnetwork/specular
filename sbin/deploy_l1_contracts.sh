#!/bin/bash

set -e

# Set workspace directory variable using pwd command
WORKSPACE_DIR=$(pwd)

# Set sbin directory variable using readlink and dirname commands
SBIN=$(dirname "$(readlink -f "$0")")
SBIN="$(
  cd "$SBIN"
  pwd
)"
# Source the utils.sh script from the sbin directory
. $SBIN/utils/utils.sh

# Set root directory variable
ROOT_DIR=$SBIN/..

# Check that the required dotenv files exist
require_dotenv "paths" ".paths.env"
require_dotenv "genesis" ".genesis.env"
require_dotenv "contracts" ".contracts.env"

# Set default value for auto accept variable
AUTO_ACCEPT=""

# Parse command line arguments
optspec="cy"
while getopts "$optspec" optchar; do
  case "${optchar}" in
  y)
    AUTO_ACCEPT="--yes"
    ;;
  c)
    # Remove debugging statement and set redeploy flag
    $SBIN/clean_deployment.sh
    REDEPLOY="true"
    ;;
  *)
    # Remove debugging statement and provide usage information
    echo "usage: $0 [-c][-s][-y][-h]"
    echo "-c : clean before running"
    echo "-s: generate and configure secrets"
    echo "-y : auto accept prompts"
    exit
    ;;
  esac
done

# Check if deployed file exists and handle redeployment
if test -f $WORKSPACE_DIR/.deployed; then
  if [[ ! -z ${REDEPLOY+x} ]]; then
    rm -f $WORKSPACE_DIR/.deployed
  else
    echo "Already Deployed"
    exit 0
  fi
fi

# Set contracts directory variable
CONTRACTS_DIR="$ROOT_DIR/contracts"

# Copy .contracts.env to contracts directory and confirm overwrite
confirm_overwrite $CONTRACTS_DIR/.env $AUTO_ACCEPT
cp .contracts.env $CONTRACTS_DIR/.env

# Get relative paths for certain files
BASE_ROLLUP_CFG_PATH=$(relative_path $BASE_ROLLUP_CFG_PATH $CONTRACTS_DIR)
ROLLUP_CFG_PATH=$(relative_path $ROLLUP_CFG_PATH $CONTRACTS_DIR)
GENESIS_PATH=$(relative_path $GENESIS_PATH $CONTRACTS_DIR)
GENESIS_CFG_PATH=$(relative_path $GENESIS_CFG_PATH $CONTRACTS_DIR)
GENESIS_EXPORTED_HASH_PATH=$(relative_path $GENESIS_EXPORTED_HASH_PATH $CONTRACTS_DIR)
DEPLOYMENTS_CFG_PATH=$(relative_path ".deployments.env" $CONTRACTS_DIR)

# Deploy contracts
cd $CONTRACTS_DIR
echo "Deploying l1 contracts..."
echo $GENESIS_EXPORTED_HASH_PATH
npx $AUTO_ACCEPT hardhat deploy --network $L1_NETWORK

# Generate deployments config and confirm overwrite
confirm_overwrite $DEPLOYMENTS_CFG_PATH $AUTO_ACCEPT
npx ts-node scripts/config/create_deployments_config.ts \
  --deployments $CONTRACTS_DIR/deployments/$L1_NETWORK \
  --deployments-config-path $DEPLOYMENTS_CFG_PATH

# Generate genesis file
cd $WORKSPACE_DIR
$SBIN/create_genesis.sh

# Generate rollup config and confirm overwrite
cd $CONTRACTS_DIR
echo "Generating rollup config..."
confirm_overwrite $ROLLUP_CFG_PATH
npx $AUTO_ACCEPT ts-node scripts/config/create_config.ts \
  --in $BASE_ROLLUP_CFG_PATH \
  --out $ROLLUP_CFG_PATH \
  --deployments $CONTRACTS_DIR/deployments/$L1_NETWORK \
  --deployments-config-path $DEPLOYMENTS_CFG_PATH \
  --genesis-path $GENESIS_PATH \
  --genesis-config-path $GENESIS_CFG_PATH \
  --genesis-hash-path $GENESIS_EXPORTED_HASH_PATH \
  --l1-network $L1_ENDPOINT

# Add deployment addresses to contracts env file
cat $DEPLOYMENTS_CFG_PATH >>$CONTRACTS_DIR/.env

# Initialize Rollup contract genesis state
echo "Initializing Rollup contract genesis state..."
npx hardhat run --network $L1_NETWORK scripts/config/set_rollup_genesis_state.ts

# Signal that deployment is done
touch $WORKSPACE_DIR/.deployed

# Print completion message
echo "Done."

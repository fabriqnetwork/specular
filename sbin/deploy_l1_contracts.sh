#!/bin/bash

WORKSPACE_DIR=$(pwd)
# the local sbin paths are relative to the project root
SBIN=$(dirname "$(readlink -f "$0")")
SBIN="$(
  cd "$SBIN"
  pwd
)"
. $SBIN/utils/utils.sh
ROOT_DIR=$SBIN/..

# Check that the all required dotenv files exists.
reqdotenv "paths" ".paths.env"
reqdotenv "genesis" ".genesis.env"
reqdotenv "contracts" ".contracts.env"

# Parse args.
optspec="c"
while getopts "$optspec" optchar; do
  case "${optchar}" in
  c)
    echo "Cleaning deployment..."
    $SBIN/clean_deployment.sh
    ;;
  *)
    echo "usage: $0 [-c][-h]"
    echo "-c : clean before running"
    exit
    ;;
  esac
done

echo "Using $CONTRACTS_DIR as HH proj"

# Copy .contracts.env
guard_overwrite $CONTRACTS_DIR/.env
cp .contracts.env $CONTRACTS_DIR/.env

# Get relative paths, since we have to run `create_genesis.sh`
# and `create_config.ts` from the HH proj.
BASE_ROLLUP_CFG_PATH=$(relpath $BASE_ROLLUP_CFG_PATH $CONTRACTS_DIR)
ROLLUP_CFG_PATH=$(relpath $ROLLUP_CFG_PATH $CONTRACTS_DIR)
GENESIS_PATH=$(relpath $GENESIS_PATH $CONTRACTS_DIR)
GENESIS_EXPORTED_HASH_PATH=$(relpath $GENESIS_EXPORTED_HASH_PATH $CONTRACTS_DIR)
DEPLOYMENTS_CFG_PATH=$(relpath ".deployments.env" $CONTRACTS_DIR)

# Deploy contracts
cd $CONTRACTS_DIR
guard "Deploy contracts? [y/N]"
echo "Deploying l1 contracts..."
npx hardhat deploy --network $L1_NETWORK

echo "Generating deployments config..."
guard_overwrite $DEPLOYMENTS_CFG_PATH
npx ts-node scripts/config/create_deployments_config.ts \
  --deployments $CONTRACTS_DIR/deployments/$L1_NETWORK \
  --deployments-config-path $DEPLOYMENTS_CFG_PATH

# Generate genesis file
cd $WORKSPACE_DIR
$SBIN/create_genesis.sh

# Generate rollup config
cd $CONTRACTS_DIR
echo "Generating rollup config..."
guard_overwrite $ROLLUP_CFG_PATH
npx ts-node scripts/config/create_config.ts \
  --in $BASE_ROLLUP_CFG_PATH \
  --out $ROLLUP_CFG_PATH \
  --deployments $CONTRACTS_DIR/deployments/$L1_NETWORK \
  --deployments-config-path $DEPLOYMENTS_CFG_PATH \
  --genesis $GENESIS_PATH \
  --genesis-hash-path $GENESIS_EXPORTED_HASH_PATH \
  --l1-network $L1_ENDPOINT

# Add deployment addresses to contracts env file
cat $DEPLOYMENTS_CFG_PATH >>$CONTRACTS_DIR/.env

echo "Initializing Rollup contract genesis state..."
npx ts-node scripts/config/set_rollup_genesis_state.ts \
  --deployments $CONTRACTS_DIR/deployments/$L1_NETWORK

echo "Done."

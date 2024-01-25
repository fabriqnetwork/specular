#!/bin/bash
set -e

WORKSPACE_DIR=$HOME/.spc/workspaces/active_workspace

# the local sbin paths are relative to the project root
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

reqdotenv "paths" $PATHS_ENV
reqdotenv "genesis" $GENESIS_ENV
reqdotenv "contracts" $CONTRACTS_ENV

AUTO_ACCEPT=""

# Parse args.
optspec="cy"
while getopts "$optspec" optchar; do
  case "${optchar}" in
  y)
    AUTO_ACCEPT="--yes"
    ;;
  c)
    echo "Cleaning deployment..."
    $SBIN/clean_deployment.sh
    REDEPLOY="true"
    ;;
  *)
    echo "usage: $0 [-c][-s][-y][-h]"
    echo "-c : clean before running"
    echo "-s: generate and configure secrets"
    echo "-y : auto accept prompts"
    exit
    ;;
  esac
done

if test -f $WORKSPACE_DIR/.deployed; then
  if [[ ! -z ${REDEPLOY+x} ]]; then
    rm -f $WORKSPACE_DIR/.deployed
  else
    echo "Already Deployed"
    exit 0
  fi
fi

echo "Using $CONTRACTS_DIR as HH proj"

# Copy .contracts.env
guard_overwrite $CONTRACTS_DIR/.env $AUTO_ACCEPT
cp $CONTRACTS_ENV $CONTRACTS_DIR/.env

# Get relative paths, since we have to run `create_genesis.sh`
# and `create_config.ts` from the HH proj.
BASE_ROLLUP_CFG_PATH=$(relpath $BASE_ROLLUP_CFG_PATH $CONTRACTS_DIR)
ROLLUP_CFG_PATH=$(relpath $ROLLUP_CFG_PATH $CONTRACTS_DIR)
GENESIS_PATH=$(relpath $GENESIS_PATH $CONTRACTS_DIR)
GENESIS_CFG_PATH=$(relpath $GENESIS_CFG_PATH $CONTRACTS_DIR)
GENESIS_EXPORTED_HASH_PATH=$(relpath $GENESIS_EXPORTED_HASH_PATH $CONTRACTS_DIR)
DEPLOYMENTS_CFG_PATH=$(relpath "$WORKSPACE_DIR/.deployments.env" $CONTRACTS_DIR)

echo $DEPLOYMENTS_CFG_PATH

# Deploy contracts
cd $CONTRACTS_DIR
# guard "Deploy contracts? [y/N]"
echo "Deploying l1 contracts..."
echo $GENESIS_EXPORTED_HASH_PATH
npx $AUTO_ACCEPT hardhat deploy --network $L1_NETWORK

echo "Generating deployments config..."
guard_overwrite $DEPLOYMENTS_CFG_PATH $AUTO_ACCEPT
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

echo "Initializing Rollup contract genesis state..."
npx hardhat run --network $L1_NETWORK scripts/config/set_rollup_genesis_state.ts

# Signal that we're done.
touch $WORKSPACE_DIR/.deployed

echo "Done."

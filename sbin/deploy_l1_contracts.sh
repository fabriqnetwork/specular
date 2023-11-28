#!/bin/bash

# the local sbin paths are relative to the project root
SBIN=$(dirname "$(readlink -f "$0")")
SBIN="`cd "$SBIN"; pwd`"
ROOT_DIR=$SBIN/..

# Check that the all required dotenv files exists.
PATHS_ENV=".paths.env"
if ! test -f "$PATHS_ENV"; then
    echo "Expected paths dotenv at $PATHS_ENV (does not exist)."
    exit
fi
echo "Using paths dotenv: $PATHS_ENV"
. $PATHS_ENV

GENESIS_ENV=".genesis.env"
if ! test -f "$GENESIS_ENV"; then
    echo "Expected dotenv at $GENESIS_ENV (does not exist)."
    exit
fi
echo "Using genesis dotenv: $GENESIS_ENV"
. $GENESIS_ENV

CONTRACTS_ENV=".contracts.env"
if  ! test -f "$CONTRACTS_ENV"; then
    echo "Expected dotenv at $CONTRACTS_ENV (does not exist)."
    exit
fi
echo "Using contracts dotenv: $CONTRACTS_ENV"
. $CONTRACTS_ENV

# Parse args.
optspec="ch"
while getopts "$optspec" optchar; do
    case "${optchar}" in
        c)
	    echo "Cleaning..."
	    $SBIN/clean_deployment.sh
	    ;;
        h)
            echo "usage: $0 [-c][-h]"
	    echo "-c : clean before running"
            exit
            ;;
        *)
            if [ "$OPTERR" != 1 ] || [ "${optspec:0:1}" = ":" ]; then
                echo "Unknown option: '-${OPTARG}'"
		exit 1
            fi
            ;;
    esac
done

echo "Using $CONTRACTS_DIR as HH proj"

# Define a function to convert a path to be relative to another directory.
relpath () {
    echo `python3 -c "import os.path; print(os.path.relpath('$1', '$2'))"`
}

# Define a function that requests a user to confirm
# that overwriting file ($1) is okay, if it exists.
guard_overwrite () {
  if test -f $1; then
    read -r -p "Overwrite $1 with a new file? [y/N] " response
    if [[ "$response" =~ ^([yY][eE][sS]|[yY])$ ]]; then
      rm $1
    else
      exit
    fi
  fi
}

# Copy .contracts.env
guard_overwrite $CONTRACTS_DIR/.env
cp $CONTRACTS_ENV $CONTRACTS_DIR/.env

# Get relative paths, since we have to run `create_genesis.ts`
# and `create_config.ts` from the HH proj.
BASE_ROLLUP_CFG_PATH=`relpath $BASE_ROLLUP_CFG_PATH $CONTRACTS_DIR`
ROLLUP_CFG_PATH=`relpath $ROLLUP_CFG_PATH $CONTRACTS_DIR`
GENESIS_PATH=`relpath $GENESIS_PATH $CONTRACTS_DIR`
GENESIS_EXPORTED_HASH_PATH=`relpath $GENESIS_EXPORTED_HASH_PATH $CONTRACTS_DIR`

# Generate genesis file
$SBIN/create_genesis.sh

# Deploy contracts
cd $CONTRACTS_DIR
echo "Deploying l1 contracts..."
echo $GENESIS_EXPORTED_HASH_PATH
npx hardhat deploy --network $L1_NETWORK

# Generate rollup config
echo "Generating rollup config..."
guard_overwrite $ROLLUP_CFG_PATH
npx ts-node scripts/config/create_config.ts \
  --in $BASE_ROLLUP_CFG_PATH \
  --out $ROLLUP_CFG_PATH \
  --genesis $GENESIS_PATH \
  --genesis-hash-path $GENESIS_EXPORTED_HASH_PATH \
  --l1-network $L1_ENDPOINT

echo "Done."

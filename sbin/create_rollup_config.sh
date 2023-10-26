#!/bin/bash
if [ ! -d "$CONTRACTS_DIR" ]; then
    SBIN=`dirname $0`
    SBIN="`cd "$SBIN"; pwd`"
    . $SBIN/configure.sh
    CONTRACTS_DIR="`cd "$CONTRACTS_DIR"; pwd`"
fi
echo "Using $CONTRACTS_DIR as HH proj"

# Check that the dotenv exists.
ENV=".genesis.env"
if ! test -f $ENV; then
    echo "Expected dotenv at $ENV (does not exist)."
    exit
fi
echo "Using dotenv: $ENV"
. $ENV

# Define a function to convert a path to be relative to another directory.
relpath () {
    echo `python3 -c "import os.path; print(os.path.relpath('$1', '$2'))"`
}

# Get relative paths, since we have to run `create_genesis.ts` from the HH proj.
BASE_ROLLUP_CFG_PATH=`relpath $BASE_ROLLUP_CFG_PATH $CONTRACTS_DIR`
ROLLUP_CFG_PATH=`relpath $ROLLUP_CFG_PATH $CONTRACTS_DIR`

# Create genesis.json file.
cd $CONTRACTS_DIR && npx ts-node scripts/config/create_config.ts \
  --in $BASE_ROLLUP_CFG_PATH \
  --out $ROLLUP_CFG_PATH

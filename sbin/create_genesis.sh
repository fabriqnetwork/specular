#!/bin/bash
if [ -d $CONFIG_DIR ]; then
    SBIN=`dirname $0`
    SBIN="`cd "$SBIN"; pwd`"
    . $SBIN/configure.sh
    CONFIG_DIR="`cd "$CONFIG_DIR"; pwd`"
fi
echo "Using $CONFIG_DIR as HH proj"

# Check that the dotenv exists.
ENV=".genesis.env"
if ! test -f $ENV; then
    echo "Expected dotenv at $ENV (does not exist)."
    exit
fi
echo "Using dotenv: $ENV"
. $ENV

relpath () {
    echo `python3 -c "import os.path; print(os.path.relpath('$1', '$2'))"`
}

# Get relative paths, since we have to run `create_genesis.ts` from the HH proj.
BASE_GENESIS_PATH=`relpath $BASE_GENESIS_PATH $CONFIG_DIR`
GENESIS_PATH=`relpath $GENESIS_PATH $CONFIG_DIR`

# Create genesis.json file.
cd $CONFIG_DIR
npx ts-node src/create_genesis.ts --in $BASE_GENESIS_PATH --out $GENESIS_PATH

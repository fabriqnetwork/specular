#!/bin/bash
if [ -d $CONFIG_DIR ]; then
    SBIN=`dirname $0`
    SBIN="`cd "$SBIN"; pwd`"
    . $SBIN/configure.sh
    CONFIG_DIR="`cd "$CONFIG_DIR"; pwd`"
fi
echo "Using $CONFIG_DIR as HH proj"

# Check that the dotenv exists.
ENV=".sp_geth.env"
if ! test -f $ENV; then
    echo "Expected dotenv at $ENV (does not exist)."
    exit
fi
echo "Using dotenv: $ENV"
. $ENV

# Get relative paths, since we have to run `create_genesis.ts` from the HH proj.
# TODO: get rid of this hack
BASE_GENESIS_PATH=`python3 -c "import os.path; print(os.path.relpath('$BASE_GENESIS_PATH', '$CONFIG_DIR'))"`
GENESIS_PATH=`python3 -c "import os.path; print(os.path.relpath('$GENESIS_PATH', '$CONFIG_DIR'))"`

cd $CONFIG_DIR
npx ts-node src/create_genesis.ts --in $BASE_GENESIS_PATH --out $GENESIS_PATH

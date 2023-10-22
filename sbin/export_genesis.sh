#!/bin/bash
# Note: do not add any printing to stdout in the positive case.
if [ -z $SP_GETH ]; then
    # If no binary specified, assume repo directory structure.
    SBIN=`dirname $0`
    SBIN="`cd "$SBIN"; pwd`"
    . $SBIN/configure.sh
    SP_GETH=$GETH_BIN
fi

# Check that the dotenv exists.
ENV=".genesis.env"
if ! test -f $ENV; then
    echo "Expected dotenv at $ENV (does not exist)."
    exit
fi
. $ENV

# Export l2 genesis hash for $GENESIS_PATH
DATA_DIR=tmp_data/
$SP_GETH init --datadir $DATA_DIR $GENESIS_PATH
$SP_GETH --datadir $DATA_DIR --verbosity 0 dump 0
rm -r $DATA_DIR

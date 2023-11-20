#!/bin/bash
SBIN=$(dirname "$(readlink -f "$0")")
ROOT_DIR=$SBIN/..

CONFIGURE_ENV=".configure.env"
if ! test -f $CONFIGURE_ENV; then
    echo "Expected dotenv at $CONFIGURE_ENV (does not exist)."
    exit
fi
echo "Using dotenv: $CONFIGURE_ENV"
. $CONFIGURE_ENV

GENESIS_ENV=".genesis.env"
if test -f $GENESIS_ENV; then
    . $GENESIS_ENV
fi

if test -f "$GENESIS_PATH"; then
    echo "Removing $GENESIS_PATH"
    rm $GENESIS_PATH
fi
if test -f "$GENESIS_EXPORTED_HASH_PATH"; then
    echo "Removing $GENESIS_EXPORTED_HASH_PATH"
    rm $GENESIS_EXPORTED_HASH_PATH
fi
if test -f "$ROLLUP_CFG_PATH"; then
    echo "Removing $ROLLUP_CFG_PATH"
    rm $ROLLUP_CFG_PATH
fi

echo "Removing deployment files in $CONTRACTS_DIR"
rm -rf $CONTRACTS_DIR/deployments/*

#!/bin/bash
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

if [ ! -d "$CONTRACTS_DIR" ]; then
    SBIN=`dirname $0`
    SBIN="`cd "$SBIN"; pwd`"
    . $SBIN/configure.sh
    CONTRACTS_DIR="`cd "$CONTRACTS_DIR"; pwd`"
fi
echo "Removing deployment files in $CONTRACTS_DIR"
rm -rf $CONTRACTS_DIR/deployments/*

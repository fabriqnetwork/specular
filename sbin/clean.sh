#!/bin/bash
SBIN=`dirname $0`
$SBIN/clean_sp_geth.sh
$SBIN/clean_deployment.sh

GENESIS_ENV=".genesis.env"
if test -f "$GENESIS_ENV"; then
    . $GENESIS_ENV
    echo "Removing base config files..."
    if test -f "$BASE_GENESIS_PATH"; then
	echo "Removing $BASE_GENESIS_PATH"
	rm $BASE_GENESIS_PATH
    fi
    if test -f "$BASE_ROLLUP_CFG_PATH"; then
	echo "Removing $ROLLUP_CFG_PATH"
	rm $ROLLUP_CFG_PATH
    fi
fi

echo "Removing dotenv files..."
rm -f .contracts.env
rm -f .genesis.env
rm -f .sp_geth.env
rm -f .sp_magi.env
rm -f .sidecar.env

# From .sp_geth.env
GETH_ENV=".sp_geth.env"
if test -f "$GETH_ENV"; then
    . $GETH_ENV
    echo "Removing $JWT_SECRET_PATH"
    rm $JWT_SECRET_PATH
fi
echo "Done."

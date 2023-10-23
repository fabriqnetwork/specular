#!/bin/bash
# Check that a dotenv exists.
GETH_ENV=".sp_geth.env"
if ! test -f $GETH_ENV; then
    echo "expected dotenv at ./$GETH_ENV (does not exist); could not clean cwd."
    exit
fi
echo "Cleaning deployment for dotenv: $GETH_ENV"
. $GETH_ENV

GENESIS_ENV=".genesis.env"
if test -f $GENESIS_ENV; then
    . $GENESIS_ENV
fi

echo "Removing dotenv files..."
rm .genesis.env
rm .sp_geth.env
rm .sp_magi.env
rm .sidecar.env
echo "Removing associated json configs..."
# From .sp_geth.env
if test -f $JWT_SECRET_PATH; then
    echo "Removing $JWT_SECRET_PATH"
    rm $JWT_SECRET_PATH
fi
# From .genesis.env
if test -f $BASE_GENESIS_PATH; then
    echo "Removing $BASE_GENESIS_PATH"
    rm $BASE_GENESIS_PATH
fi
if test -f $GENESIS_PATH; then
    echo "Removing $GENESIS_PATH"
    rm $GENESIS_PATH
fi
echo "Removing data dir..."
rm -rf $DATA_DIR

echo "Removing contract deployment files..."
rm -rf $CONTRACTS_DIR/deployments/localhost
rm -rf $CONTRACTS_DIR/deployments/specularLocalDev

# Remove L1 docker container
# TODO: does this belong here?
echo "Force-removing docker container if it exists..."
docker rm --force l1_geth

echo "Done."

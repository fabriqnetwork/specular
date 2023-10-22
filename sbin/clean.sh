#!/bin/bash
# Check that the dotenv exists.
ENV=".sp_geth.env"
if ! test -f $ENV; then
    echo "Expected dotenv at $ENV (does not exist)."
    exit
fi
echo "Cleaning deployment for dotenv: $ENV"
. $ENV

echo "Removing env files..."
rm .sp_geth.env
rm .sp_magi.env
rm .sidecar.env
echo "Removing data dir..."
rm -rf $DATA_DIR

echo "Removing contract deployment files..."
rm -rf $CONTRACTS_DIR/deployments/localhost
rm -rf $CONTRACTS_DIR/deployments/specularLocalDev

# Remove L1 docker container
# TODO: does this belong here?
echo "Removing docker container..."
docker rm --force l1_geth

echo "Done."

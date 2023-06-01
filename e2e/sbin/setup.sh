#!/bin/bash

# Define directory structure
SBIN_DIR=`dirname $0`
SBIN_DIR="`cd "$SBIN_DIR"; pwd`"
PROJECT_DIR=$SBIN_DIR/../project
PROJECT_DATA_DIR=$PROJECT_DIR/specular-datadir
CONFIG_DIR=$SBIN_DIR/../../config
CONTRACTS_DIR=$SBIN_DIR/../../contracts
GETH_SPECULAR_DIR=$SBIN_DIR/../../clients/geth/specular

# Make project directory
rm -rf $PROJECT_DIR
mkdir -p $PROJECT_DATA_DIR

# Add keys
cp $GETH_SPECULAR_DIR/data/keys/sequencer.prv $PROJECT_DATA_DIR/key.prv
cp $GETH_SPECULAR_DIR/data/password.txt $PROJECT_DATA_DIR/password.txt

# Build and add genesis.json
cd $CONTRACTS_DIR
npx hardhat compile
cd $CONFIG_DIR
npx ts-node src/create_genesis.ts --in data/base_genesis.json --out $GETH_SPECULAR_DIR/data/genesis.json
cp $GETH_SPECULAR_DIR/data/genesis.json $PROJECT_DATA_DIR/genesis.json


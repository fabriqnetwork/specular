#!/bin/bash

# Configure variables
SBIN_DIR=`dirname $0`
SBIN_DIR="`cd "$SBIN_DIR"; pwd`"
set -o allexport
source $SBIN_DIR/configure.sh
set +o allexport

# Make project directory
rm -rf $PROJECT_DIR
mkdir -p $PROJECT_LOG_DIR
mkdir -p $PROJECT_DATA_DIR

# Add keys
cp $DATA_DIR/sequencer.prv $PROJECT_DATA_DIR/sequencer.prv
cp $DATA_DIR/validator.prv $PROJECT_DATA_DIR/validator.prv
cp $DATA_DIR/password.txt $PROJECT_DATA_DIR/password.txt

# Build and add genesis.json
cd $CONTRACTS_DIR
npm install ganache --global
npx hardhat compile
cd $CONFIG_DIR
npx ts-node src/create_genesis.ts --in data/base_genesis.json --out $GETH_SPECULAR_DIR/data/genesis.json
cp $GETH_SPECULAR_DIR/data/genesis.json $PROJECT_DATA_DIR/genesis.json

# Build L2 client
cd $GETH_SPECULAR_DIR
make geth

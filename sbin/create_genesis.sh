#!/bin/bash
SBIN=`dirname $0`
SBIN="`cd "$SBIN"; pwd`"
. $SBIN/configure.sh

cd $CONFIG_DIR

npx ts-node src/create_genesis.ts --in data/base_genesis.json --out ../e2e/data/genesis.json


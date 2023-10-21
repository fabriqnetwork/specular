#!/bin/bash
SBIN=`dirname $0`
SBIN="`cd "$SBIN"; pwd`"
. $SBIN/configure.sh
cd $CONFIG_DIR
npx ts-node src/create_genesis.ts --in $BASE_GENESIS_PATH --out $GENESIS_PATH

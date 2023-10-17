#!/bin/bash
SBIN=`dirname $0`
SBIN="`cd "$SBIN"; pwd`"
. $SBIN/configure.sh

cd $DATA_DIR

$GETH_BIN --password ./password.txt --datadir . account import ./sequencer.prv
$GETH_BIN --password ./password.txt --datadir . account import ./validator.prv

$GETH_BIN --datadir . --networkid $NETWORK_ID init ./genesis.json

#!/bin/bash
SBIN=`dirname $0`
SBIN="`cd "$SBIN"; pwd`"
. $SBIN/configure.sh

cd $DATA_DIR

$RIPCORD_BIN --password ./password.txt --datadir . account import ./sequencer.prv
$RIPCORD_BIN --password ./password.txt --datadir . account import ./validator.prv

$RIPCORD_BIN --datadir . --networkid $NETWORK_ID init ./genesis.json

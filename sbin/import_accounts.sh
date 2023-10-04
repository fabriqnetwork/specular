#!/bin/bash
SBIN=`dirname $0`
SBIN="`cd "$SBIN"; pwd`"
. $SBIN/configure.sh

cd $DATA_DIR

$SIDECAR_DIR/build/bin/geth --password ./password.txt --datadir ./data_sequencer account import ./sequencer.prv
$SIDECAR_DIR/build/bin/geth --password ./password.txt --datadir ./data_sequencer account import ./validator.prv
$SIDECAR_DIR/build/bin/geth --password ./password.txt --datadir ./data_validator account import ./sequencer.prv
$SIDECAR_DIR/build/bin/geth --password ./password.txt --datadir ./data_validator account import ./validator.prv
$SIDECAR_DIR/build/bin/geth --password ./password.txt --datadir ./data_indexer account import ./sequencer.prv
$SIDECAR_DIR/build/bin/geth --password ./password.txt --datadir ./data_indexer account import ./validator.prv

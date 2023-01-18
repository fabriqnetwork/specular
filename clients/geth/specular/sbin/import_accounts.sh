#!/bin/bash
SBIN=`dirname $0`
SBIN="`cd "$SBIN"; pwd`"
. $SBIN/configure.sh
cd $DATA_DIR
$GETH_SPECULAR_DIR/build/bin/geth --password ./password.txt --datadir ./data_sequencer account import ./keys/sequencer.prv
$GETH_SPECULAR_DIR/build/bin/geth --password ./password.txt --datadir ./data_sequencer account import ./keys/validator.prv
$GETH_SPECULAR_DIR/build/bin/geth --password ./password.txt --datadir ./data_validator account import ./keys/sequencer.prv
$GETH_SPECULAR_DIR/build/bin/geth --password ./password.txt --datadir ./data_validator account import ./keys/validator.prv
$GETH_SPECULAR_DIR/build/bin/geth --password ./password.txt --datadir ./data_indexer account import ./keys/sequencer.prv
$GETH_SPECULAR_DIR/build/bin/geth --password ./password.txt --datadir ./data_indexer account import ./keys/validator.prv

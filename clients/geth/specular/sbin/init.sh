#!/bin/bash
SBIN=`dirname $0`
SBIN="`cd "$SBIN"; pwd`"
. $SBIN/configure.sh
cd $DATA_DIR
$GETH_SPECULAR_DIR/build/bin/geth --datadir ./data_sequencer --networkid 13527 init ./genesis.json
$GETH_SPECULAR_DIR/build/bin/geth --datadir ./data_validator --networkid 13527 init ./genesis.json
$GETH_SPECULAR_DIR/build/bin/geth --datadir ./data_indexer --networkid 13527 init ./genesis.json

cd $SBIN
./init_clef.exp


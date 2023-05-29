#!/bin/bash
SBIN=`dirname $0`
SBIN="`cd "$SBIN"; pwd`"
. $SBIN/configure.sh
. $SBIN/configure_system.sh
cd $DATA_DIR

args=(
    --datadir ./data_sequencer
    --verbosity 0
)
$GETH_SPECULAR_DIR/build/bin/geth "${args[@]}" dump 0

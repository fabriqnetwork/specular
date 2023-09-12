#!/bin/bash
SBIN=`dirname $0`
SBIN="`cd "$SBIN"; pwd`"
. $SBIN/configure.sh
. $SBIN/configure_system.sh
cd $DATA_DIR

mkdir tmp_data/

args=(
    --datadir ./tmp_data/
    --verbosity 0
)

$GETH_SPECULAR_DIR/build/bin/geth init --datadir ./tmp_data ./genesis.json
$GETH_SPECULAR_DIR/build/bin/geth "${args[@]}" dump 0

rm -r tmp_data/

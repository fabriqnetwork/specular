#!/bin/bash
SBIN=`dirname $0`
SBIN="`cd "$SBIN"; pwd`"
. $SBIN/configure.sh

cd $DATA_DIR

args=(
    --datadir ./tmp_data/
    --verbosity 0
)

mkdir tmp_data/

$GETH_DIR/build/bin/geth init --datadir ./tmp_data ./genesis.json
$GETH_DIR/build/bin/geth "${args[@]}" dump 0

rm -r tmp_data/

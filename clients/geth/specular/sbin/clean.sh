#!/bin/bash
SBIN=`dirname $0`
SBIN="`cd "$SBIN"; pwd`"
. $SBIN/configure.sh
cd $DATA_DIR
rm -rf ./data_sequencer/geth
rm -rf ./data_validator/geth
rm -rf ./data_indexer/geth

#!/bin/bash
SBIN=`dirname $0`
SBIN="`cd "$SBIN"; pwd`"
. $SBIN/configure.sh
cd $DATA_DIR
../build/bin/geth --datadir ./data_sequencer attach http://localhost:4011

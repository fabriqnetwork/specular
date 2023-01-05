#!/bin/bash
SBIN=`dirname $0`
SBIN="`cd "$SBIN"; pwd`"
. $SBIN/configure.sh
cd $DATA_DIR
../build/bin/geth --password ./password.txt --datadir ./data account import ./keys/sequencer.prv
../build/bin/geth --password ./password.txt --datadir ./data account import ./keys/validator.prv
../build/bin/geth --password ./password.txt --datadir ./data_validator account import ./keys/sequencer.prv
../build/bin/geth --password ./password.txt --datadir ./data_validator account import ./keys/validator.prv
../build/bin/geth --password ./password.txt --datadir ./data_indexer account import ./keys/sequencer.prv
../build/bin/geth --password ./password.txt --datadir ./data_indexer account import ./keys/validator.prv

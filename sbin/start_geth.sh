#!/bin/bash
SBIN=`dirname $0`
SBIN="`cd "$SBIN"; pwd`"
. $SBIN/configure.sh

cd $DATA_DIR

args=(
    --datadir .
    --http --http.addr '0.0.0.0' --http.port 4011 --http.api 'engine,personal,eth,net,web3,txpool,miner,debug'
    --ws --ws.addr '0.0.0.0' --ws.port 4012 --ws.api 'engine,personal,eth,net,web3,txpool,miner,debug'
    --http.corsdomain '*' --ws.origins '*'
    --http.vhosts '*'
    --networkid $NETWORK_ID
)

$GETH_BIN "${args[@]}"

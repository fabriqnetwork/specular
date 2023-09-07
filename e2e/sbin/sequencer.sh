#!/bin/bash

# Configure variables
SBIN_DIR=`dirname $0`
SBIN_DIR="`cd "$SBIN_DIR"; pwd`"
set -o allexport
source $SBIN_DIR/configure.sh
set +o allexport

# Import account
$L2GETH_BIN --datadir . --password ./password.txt account import ./sequencer.prv
$L2GETH_BIN --datadir . --password ./password.txt account import ./validator.prv

# Initialize geth
$L2GETH_BIN --datadir . --networkid $NETWORK_ID init ./genesis.json

# Run geth
exec $L2GETH_BIN \
    --password ./password.txt \
    --datadir . \
    --networkid $NETWORK_ID \
    --nodiscover \
    --maxpeers 0 \
    --http \
    --http.port=$L2_HTTP_PORT \
    --http.addr=0.0.0.0 \
    --http.corsdomain=* \
    --http.api=personal,eth,net,web3,txpool,miner,proof,debug \
    --ws \
    --ws.port=$L2_WS_PORT \
    --ws.addr=0.0.0.0 \
    --ws.origins=* \
    --ws.api=personal,eth,net,web3,txpool,miner,proof,debug \
    --rollup.node=sequencer \
    --rollup.coinbase=$SEQUENCER_ADDR \
    --rollup.l1endpoint=$L1_ENDPOINT \
    --rollup.l1chainid=$L1_CHAIN_ID \
    --rollup.sequencer-inbox-addr=$SEQUENCER_INBOX_ADDR \
    --rollup.rollup-addr=$ROLLUP_ADDR \
    --rollup.rollup-stake-amount=$ROLLUP_STAKE_AMOUNT

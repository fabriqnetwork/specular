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
args=(
    --datadir .
    --password ./password.txt
    --http --http.addr '0.0.0.0' --http.port 4011 --http.api 'personal,eth,net,web3,txpool,miner,proof,debug'
    --ws --ws.addr '0.0.0.0' --ws.port 4012 --ws.api 'personal,eth,net,web3,txpool,miner,proof,debug'
    --http.corsdomain '*' --ws.origins '*'
    --networkid $NETWORK_ID
    --rollup.l1.endpoint $L1_ENDPOINT
    --rollup.l1.chainid $L1_CHAIN_ID
    --rollup.l1.sequencer-inbox-addr $SEQUENCER_INBOX_ADDR
    --rollup.l1.rollup-addr $ROLLUP_ADDR
    --rollup.l1.rollup-genesis-block $GENESIS_L1_BLOCK_NUM
    --rollup.l2.chainid $NETWORK_ID
    --rollup.sequencer.addr $SEQUENCER_ADDR
    --rollup.validator.addr $VALIDATOR_ADDR
)

exec $L2GETH_BIN "${args[@]}"

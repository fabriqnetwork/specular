#!/bin/bash
SBIN=`dirname $0`
SBIN="`cd "$SBIN"; pwd`"
. $SBIN/configure.sh
cd $DATA_DIR

$GETH_SPECULAR_DIR/build/bin/geth \
    --datadir ./data_sequencer \
    --password ./password.txt \
    --http --http.addr '0.0.0.0' --http.port 4011 --http.api 'personal,eth,net,web3,txpool,miner,proof,debug' \
    --ws --ws.addr '0.0.0.0' --ws.port 4012 --ws.api 'personal,eth,net,web3,txpool,miner,proof,debug' \
    --http.corsdomain '*' --ws.origins '*' \
    --networkid 13527 \
    --rollup.node 'sequencer' \
    --rollup.coinbase=f39fd6e51aad88f6f4ce6ab8827279cfffb92266 \
    --rollup.l1endpoint 'ws://localhost:8545' \
    --rollup.l1chainid 31337 \
    --rollup.sequencer-inbox-addr '0x2E983A1Ba5e8b38AAAeC4B440B9dDcFBf72E15d1' \
    --rollup.rollup-addr '0xF6168876932289D073567f347121A267095f3DD6' \
    --rollup.rollup-stake-amount 100 \
    --maxpeers 0

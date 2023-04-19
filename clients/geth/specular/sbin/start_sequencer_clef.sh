#!/bin/bash
SBIN=`dirname $0`
SBIN="`cd "$SBIN"; pwd`"
. $SBIN/configure.sh
cd $DATA_DIR

$GETH_SPECULAR_DIR/build/bin/geth \
    --datadir ./data_sequencer \
    --http --http.addr '0.0.0.0' --http.port 4011 --http.api 'personal,eth,net,web3,txpool,miner,proof,debug' \
    --ws --ws.addr '0.0.0.0' --ws.port 4012 --ws.api 'personal,eth,net,web3,txpool,miner,proof,debug' \
    --http.corsdomain '*' --ws.origins '*' \
    --networkid 13527 \
    --rollup.node 'sequencer' \
    --rollup.coinbase f39fd6e51aad88f6f4ce6ab8827279cfffb92266 \
    --rollup.clefendpoint 'http://127.0.0.1:8550/' \
    --rollup.l1endpoint 'ws://localhost:8545' \
    --rollup.l1chainid 31337 \
    --rollup.sequencer-inbox-addr '0xe7f1725E7734CE288F8367e1Bb143E90bb3F0512' \
    --rollup.rollup-addr '0x5FC8d32690cc91D4c39d9d3abcBD16989F875707' \
    --rollup.rollup-stake-amount 100

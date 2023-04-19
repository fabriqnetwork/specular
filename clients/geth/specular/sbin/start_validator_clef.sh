#!/bin/bash
SBIN=`dirname $0`
SBIN="`cd "$SBIN"; pwd`"
. $SBIN/configure.sh
cd $DATA_DIR
$GETH_SPECULAR_DIR/build/bin/geth \
    --datadir ./data_validator \
    --http --http.addr '0.0.0.0' --http.port 4018 --http.api 'personal,eth,net,web3,txpool,miner,proof,debug' \
    --ws --ws.addr '0.0.0.0' --ws.port 4019 --ws.api 'personal,eth,net,web3,txpool,miner,proof,debug' \
    --http.corsdomain '*' --ws.origins '*' \
    --networkid 13527 \
    --port 30304 \
    --authrpc.port 8561 \
    --rollup.node 'validator' \
    --rollup.coinbase=70997970c51812dc3a010c7d01b50e0d17dc79c8 \
    --rollup.clefendpoint 'http://127.0.0.1:8550/' \
    --rollup.l1endpoint 'ws://localhost:8545' \
    --rollup.l1chainid 31337 \
    --rollup.sequencer-addr '0xf39fd6e51aad88f6f4ce6ab8827279cfffb92266' \
    --rollup.sequencer-inbox-addr '0xe7f1725E7734CE288F8367e1Bb143E90bb3F0512' \
    --rollup.rollup-addr '0x5FC8d32690cc91D4c39d9d3abcBD16989F875707' \
    --rollup.rollup-stake-amount 100 \
    --rollup.l1-rollup-genesis-block 0

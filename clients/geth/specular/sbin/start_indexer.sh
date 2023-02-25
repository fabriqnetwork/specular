#!/bin/bash
SBIN=`dirname $0`
SBIN="`cd "$SBIN"; pwd`"
. $SBIN/configure.sh
cd $DATA_DIR
$GETH_SPECULAR_DIR/build/bin/geth \
    --datadir ./data_indexer \
    --password ./password.txt \
    --http --http.addr '0.0.0.0' --http.port 4021 --http.api 'personal,eth,net,web3,txpool,miner,proof,debug' \
    --ws --ws.addr '0.0.0.0' --ws.port 4022 --ws.api 'personal,eth,net,web3,txpool,miner,proof,debug' \
    --http.corsdomain '*' --ws.origins '*' \
    --gcmode=archive \
    --networkid 13527 \
    --port 30305 \
    --authrpc.port 8562 \
    --rollup.node 'indexer' \
    --rollup.coinbase=f39fd6e51aad88f6f4ce6ab8827279cfffb92266 \
    --rollup.l1endpoint 'ws://localhost:8545' \
    --rollup.l1chainid 31337 \
    --rollup.sequencer-inbox-addr '0xe7f1725E7734CE288F8367e1Bb143E90bb3F0512' \
    --rollup.rollup-addr '0x5FC8d32690cc91D4c39d9d3abcBD16989F875707' \
    --rollup.rollup-stake-amount 100 \
    --rollup.l1-rollup-genesis-block 0

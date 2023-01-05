#!/bin/bash
SBIN=`dirname $0`
SBIN="`cd "$SBIN"; pwd`"
. $SBIN/configure.sh
cd $DATA_DIR
../build/bin/geth \
    --datadir ./data_validator \
    --password ./password.txt \
    --http --http.addr '0.0.0.0' --http.port 4018 --http.api 'personal,eth,net,web3,txpool,miner,proof,debug' \
    --ws --ws.addr '0.0.0.0' --ws.port 4019 --ws.api 'personal,eth,net,web3,txpool,miner,proof,debug' \
    --http.corsdomain '*' --ws.origins '*' \
    --networkid 13527 \
    --port 30304 \
    --authrpc.port 8561 \
    --rollup.node 'validator' \
    --rollup.coinbase=70997970c51812dc3a010c7d01b50e0d17dc79c8 \
    --rollup.l1endpoint 'ws://localhost:8545' \
    --rollup.l1chainid 31337 \
    --rollup.sequencer-addr '0xf39fd6e51aad88f6f4ce6ab8827279cfffb92266' \
    --rollup.sequencer-inbox-addr '0x9fE46736679d2D9a65F0992F2272dE9f3c7fa6e0' \
    --rollup.rollup-addr '0x0165878A594ca255338adfa4d48449f69242Eb8F' \
    --rollup.rollup-stake-amount 100 \
    --rollup.l1-rollup-genesis-block 0

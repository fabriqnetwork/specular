#!/bin/bash
../../build/bin/geth \
    --datadir ../../data \
    --password password.txt \
    --http --http.addr '0.0.0.0' --http.port 4011 --http.api 'personal,eth,net,web3,txpool,miner,proof,debug' \
    --ws --ws.addr '0.0.0.0' --ws.port 4012 --ws.api 'personal,eth,net,web3,txpool,miner,proof,debug' \
    --http.corsdomain '*' --ws.origins '*' \
    --networkid 13527 \
    --rollup.node 'sequencer' \
    --rollup.coinbase=f39fd6e51aad88f6f4ce6ab8827279cfffb92266 \
    --rollup.l1endpoint 'ws://localhost:8545' \
    --rollup.l1chainid 31337 \
    --rollup.sequencer-inbox-addr '0x9fE46736679d2D9a65F0992F2272dE9f3c7fa6e0' \
    --rollup.rollup-addr '0x0165878A594ca255338adfa4d48449f69242Eb8F' \
    --rollup.rollup-stake-amount 100

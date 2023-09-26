#!/bin/bash
SBIN=`dirname $0`
SBIN="`cd "$SBIN"; pwd`"
. $SBIN/configure.sh
. $SBIN/configure_system.sh
cd $DATA_DIR

args=(
    --datadir ./data_indexer
    --password ./password.txt
    --http --http.addr '0.0.0.0' --http.port 4021 --http.api 'personal,eth,net,web3,txpool,miner,proof,debug'
    --ws --ws.addr '0.0.0.0' --ws.port 4022 --ws.api 'personal,eth,net,web3,txpool,miner,proof,debug'
    --http.corsdomain '*' --ws.origins '*'
    --http.vhosts '*'
    --gcmode=archive
    --networkid $NETWORK_ID
    --port 30305
    --authrpc.port 8562
    --rollup.l1endpoint $L1_ENDPOINT
    --rollup.l1chainid $L1_CHAIN_ID
    --rollup.l1.sequencer-inbox-addr $SEQUENCER_INBOX_ADDR
    --rollup.l1.rollup-addr $ROLLUP_ADDR
    --rollup.l1.rollup-genesis-block $GENESIS_L1_BLOCK_NUM
    --rollup.l2.chainid $NETWORK_ID
)

if $USE_CLEF == 'true'; then
    args+=(--rollup.clefendpoint $CLEF_ENDPOINT)
fi

$SIDECAR_DIR/build/bin/geth "${args[@]}"

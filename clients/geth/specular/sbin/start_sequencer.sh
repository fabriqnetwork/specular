#!/bin/bash
SBIN=`dirname $0`
SBIN="`cd "$SBIN"; pwd`"
. $SBIN/configure.sh
cd $DATA_DIR

args=(
    --datadir ./data_sequencer
    --password ./password.txt
    --http --http.addr '0.0.0.0' --http.port 4011 --http.api 'personal,eth,net,web3,txpool,miner,proof,debug'
    --ws --ws.addr '0.0.0.0' --ws.port 4012 --ws.api 'personal,eth,net,web3,txpool,miner,proof,debug'
    --http.corsdomain '*' --ws.origins '*'
    --networkid $NETWORK_ID
    --rollup.node 'sequencer'
    --rollup.coinbase $COINBASE_ADDR
    --rollup.l1endpoint $L1_ENDPOINT
    --rollup.l1chainid $L1_CHAIN_ID
    --rollup.sequencer-inbox-addr $SEQUENCER_INBOX_ADDR
    --rollup.rollup-addr $ROLLUP_ADDR
    --rollup.rollup-stake-amount $ROLLUP_STAKE_AMOUNT
)

if $USE_CLEF == 'true'; then
    args+=(--rollup.clefendpoint $CLEF_ENDPOINT)
fi

$GETH_SPECULAR_DIR/build/bin/geth "${args[@]}"

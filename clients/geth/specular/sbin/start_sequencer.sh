#!/bin/bash
SBIN=`dirname $0`
SBIN="`cd "$SBIN"; pwd`"
. $SBIN/configure.sh
. $SBIN/configure_system.sh
cd $DATA_DIR

args=(
    --datadir ./data_sequencer
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

if [[ $USE_CLEF == 'true' ]]; then
    args+=(--rollup.sequencer.clef-endpoint $CLEF_ENDPOINT)
fi

$GETH_SPECULAR_DIR/build/bin/geth "${args[@]}"

#!/bin/bash
SBIN=`dirname $0`
SBIN="`cd "$SBIN"; pwd`"
. $SBIN/configure.sh
. $SBIN/configure_system.sh
cd $DATA_DIR

args=(
    --datadir ./data_validator
    --http --http.addr '0.0.0.0' --http.port 4018 --http.api 'personal,eth,net,web3,txpool,miner,proof,debug'
    --ws --ws.addr '0.0.0.0' --ws.port 4019 --ws.api 'personal,eth,net,web3,txpool,miner,proof,debug'
    --http.corsdomain '*' --ws.origins '*'
    --networkid $NETWORK_ID
    --port 30304
    --authrpc.port 8561
    # Rollup flags
    --rollup.l1-endpoint $L1_ENDPOINT
    --rollup.l1-chainid $L1_CHAIN_ID
    --rollup.l1-rollup-genesis-block 0
    --rollup.l1-sequencer-inbox-addr $SEQUENCER_INBOX_ADDR
    --rollup.l1-rollup-addr $ROLLUP_ADDR
    --rollup.validator-addr $VALIDATOR_ADDR
    --rollup.validator-is-active-challenger
    --rollup.rollup-stake-amount $ROLLUP_STAKE_AMOUNT
)

if $USE_CLEF == 'true'; then
    args+=(--rollup.l2-clef-endpoint $CLEF_ENDPOINT)
else
    args+=(--password ./password.txt)
fi

$GETH_SPECULAR_DIR/build/bin/geth "${args[@]}"

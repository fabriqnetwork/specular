#!/bin/bash
SBIN=`dirname $0`
SBIN="`cd "$SBIN"; pwd`"
. $SBIN/configure.sh
. $SBIN/configure_system.sh
cd $DATA_DIR

args=(
    --rollup.sequencer
    --rollup.l1.endpoint $L1_ENDPOINT
    --rollup.l1.chainid $L1_CHAIN_ID
    --rollup.l1.sequencer-inbox-addr $SEQUENCER_INBOX_ADDR
    --rollup.l1.rollup-addr $ROLLUP_ADDR
    --rollup.l1.rollup-genesis-block $GENESIS_L1_BLOCK_NUM
    --rollup.l2.chainid $NETWORK_ID
    --rollup.sequencer.addr $SEQUENCER_ADDR
    --rollup.validator
    --rollup.validator.addr $VALIDATOR_ADDR
    --keystore.keystore $DATA_DIR/data_sequencer/keystore
)

if [[ $USE_CLEF == 'true' ]]; then
    args+=(--rollup.sequencer.clef-endpoint $CLEF_ENDPOINT)
fi

echo $SIDECAR_DIR/build/bin/sidecar "${args[@]}"
$SIDECAR_DIR/build/bin/sidecar "${args[@]}"

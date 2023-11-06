#!/bin/sh

# Check that the dotenv exists.
ENV=".sidecar.env"
if ! test -f $ENV; then
    echo "Expected dotenv at $ENV (does not exist)."
    exit
fi
echo "Using dotenv: $ENV"
. $ENV

if [ -z $SIDECAR ]; then
    # If no binary specified, assume repo directory structure.
    SBIN=`dirname $0`
    SBIN="`cd "$SBIN"; pwd`"
    . $SBIN/configure.sh
    SIDECAR=$SIDECAR_BIN
fi

ARGS="
    --rollup.sequencer \
    --rollup.l1.endpoint $L1_ENDPOINT \
    --rollup.l1.chainid $L1_CHAIN_ID \
    --rollup.l1.sequencer-inbox-addr $SEQUENCER_INBOX_ADDR \
    --rollup.l1.rollup-addr $ROLLUP_ADDR \
    --rollup.l1.rollup-genesis-block $GENESIS_L1_BLOCK_NUM \
    --rollup.l2.endpoint $L2_ENDPOINT \
    --rollup.l2.chainid $NETWORK_ID \
    --rollup.sequencer.addr $SEQUENCER_ADDR \
    --sequencer.secret-key $SEQUENCER_SECRET_KEY \
    --validator.secret-key $VALIDATOR_SECRET_KEY \
"

echo "starting sidecar with the following flags:"
echo $ARGS

$SIDECAR $ARGS

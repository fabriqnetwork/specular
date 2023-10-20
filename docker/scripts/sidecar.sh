#!/bin/sh

sidecar \
    --rollup.sequencer \
    --rollup.l1.endpoint $L1_ENDPOINT \
    --rollup.l1.chainid $L1_CHAIN_ID \
    --rollup.l1.sequencer-inbox-addr $SEQUENCER_INBOX_ADDR \
    --rollup.l1.rollup-addr $ROLLUP_ADDR \
    --rollup.l1.rollup-genesis-block $GENESIS_L1_BLOCK_NUM \
    --rollup.l2.chainid $NETWORK_ID \
    --rollup.sequencer.addr $SEQUENCER_ADDR \
    --rollup.validator \
    --rollup.validator.addr $VALIDATOR_ADDR \

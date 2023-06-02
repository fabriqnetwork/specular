#!/bin/sh
geth --datadir . --password ./password.txt account import ./key.prv
geth --datadir . --networkid $NETWORK_ID init ./genesis.json

DOCKER=1 IS_SEQUENCER=$IS_SEQUENCER IS_VALIDATOR=$IS_VALIDATOR IS_INDEXER=$IS_INDEXER \
    . config/sbin/configure_geth_args.sh
IS_SEQUENCER=$IS_SEQUENCER IS_VALIDATOR=$IS_VALIDATOR \
    . config/sbin/configure_specular_node_defaults.sh
# Run Geth with configured args.
exec geth \
    "$GETH_ARGS" \
    "$SPECULAR_NODE_DEFAULTS" \
    --rollup.config $SPECULAR_NODE_CFG_PATH "$@"

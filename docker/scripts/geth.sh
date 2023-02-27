#!/bin/sh

geth --datadir . --password ./password.txt account import ./key.prv

geth --datadir . --networkid $NETWORK_ID init ./genesis.json

exec geth \
    --password ./password.txt \
    --datadir . \
    --networkid $NETWORK_ID \
    --nodiscover \
    --maxpeers 0 \
    "$@"

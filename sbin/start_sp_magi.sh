#!/bin/bash

# Enforce that the dotenv exists.
ENV=".sp_magi.env"

if ! test -f $ENV; then
    echo "Expected dotenv at $ENV (does not exist)."
    exit
fi

echo "Using dotenv: $ENV"
# TODO: why does this not work
#. $ENV
. $(pwd)/$ENV

if [ -z $MAGI ]; then
    # If no binary specified, assume repo directory structure.
    SBIN=`dirname $0`
    ROOT="`cd "$SBIN/../"; pwd`"
    MAGI=$ROOT/services/cl_clients/magi/target/debug/magi
fi
echo "Using bin: $MAGI"

# Set sync flags.
SYNC_FLAGS=""
if [ $SYNC_MODE = "checkpoint" ] ; then
    SYNC_FLAGS="--checkpoint-sync-url $CHECKPOINT_SYNC_URL --checkpoint-hash $CHECKPOINT_HASH"
fi

# Set devnet flags.
DEVNET_FLAGS=""
if [ "$DEVNET" = true ] ; then
    echo "Enabling devnet mode."
    DEVNET_FLAGS="--devnet"
fi

# Set local sequencer flags.
SEQUENCER_FLAGS=""
if [ "$SEQUENCER" = true ] ; then
    echo "Enabling local sequencer."
    SEQUENCER_FLAGS="--sequencer --sequencer-max-safe-lag $SEQUENCER_MAX_SAFE_LAG"
fi

CMD="$MAGI \
    --network $NETWORK \
    --l1-rpc-url $L1_RPC_URL \
    --l2-rpc-url $L2_RPC_URL \
    --sync-mode $SYNC_MODE \
    --l2-engine-url $L2_ENGINE_URL \
    --jwt-file $JWT_SECRET_PATH \
    --rpc-port $RPC_PORT \
    $SYNC_FLAGS \
    $DEVNET_FLAGS \
    $SEQUENCER_FLAGS"

echo "$CMD"
exec $CMD

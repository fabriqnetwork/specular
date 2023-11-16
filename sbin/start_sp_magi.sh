#!/bin/bash
SBIN=`dirname $0`

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

if [ -z $SP_MAGI_BIN ]; then
    # If no binary specified, assume repo directory structure.
    . $SBIN/configure.sh
fi
echo "Using bin: $SP_MAGI_BIN"

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

FLAGS="
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

echo "starting sp-magi with the following flags:"
echo "$FLAGS"
$SP_MAGI_BIN $FLAGS
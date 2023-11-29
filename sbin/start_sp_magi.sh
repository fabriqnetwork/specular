#!/bin/bash

# the local sbin paths are relative to the project root
SBIN=$(dirname "$(readlink -f "$0")")
SBIN="`cd "$SBIN"; pwd`"
ROOT_DIR=$SBIN/..

# Check that all required dotenv files exists.
PATHS_ENV=".paths.env"
if ! test -f "$PATHS_ENV"; then
    echo "Expected dotenv at $PATHS_ENV (does not exist)."
    exit
fi
echo "Using paths dotenv: $PATHS_ENV"
. $PATHS_ENV

SP_MAGI_ENV=".sp_magi.env"
if ! test -f "$SP_MAGI_ENV"; then
    echo "Expected dotenv at $SP_MAGI_ENV (does not exist)."
    exit
fi
echo "Using sp_magi dotenv: $SP_MAGI_ENV"
. $SP_MAGI_ENV

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
    SEQUENCER_FLAGS="
        --sequencer \
        --sequencer-max-safe-lag $SEQUENCER_MAX_SAFE_LAG \
        --sequencer-pk-file $SEQUENCER_PK_FILE"
fi

# TODO: use array for flags
FLAGS="
    --network $NETWORK \
    --l1-rpc-url $L1_RPC_URL \
    --l2-rpc-url $L2_RPC_URL \
    --sync-mode $SYNC_MODE \
    --l2-engine-url $L2_ENGINE_URL \
    --jwt-file $JWT_SECRET_PATH \
    --rpc-port $RPC_PORT \
    $SYNC_FLAGS $DEVNET_FLAGS $SEQUENCER_FLAGS $@"

echo "starting sp-magi with the following flags:"
echo "$FLAGS"
$SP_MAGI_BIN $FLAGS

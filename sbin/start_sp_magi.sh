#!/bin/bash
set -e

# Set project root and sbin paths
SBIN=$(dirname "$(readlink -f "$0")")
SBIN="$(
  cd "$SBIN"
  pwd
)"
. $SBIN/utils/utils.sh
ROOT_DIR=$SBIN/..

# Check that all required dotenv files exist
require_dotenv "paths" ".paths.env"
require_dotenv "sp_magi" ".sp_magi.env"

# Generate waitfile for service init (docker/k8)
WAITFILE="/tmp/.${0##*/}.lock"
if [[ ! -z ${WAIT_DIR+x} ]]; then
  WAITFILE=$WAIT_DIR/.${0##*/}.lock
fi

# Set sync flags
SYNC_FLAGS=""
if [ $SYNC_MODE = "checkpoint" ]; then
  SYNC_FLAGS="--checkpoint-sync-url $CHECKPOINT_SYNC_URL --checkpoint-hash $CHECKPOINT_HASH"
fi

# Set devnet flags
DEVNET_FLAGS=""
if [ "$DEVNET" = true ]; then
  DEVNET_FLAGS="--devnet"
fi

# Set local sequencer flags
SEQUENCER_FLAGS=""
if [ "$SEQUENCER" = true ]; then
  SEQUENCER_FLAGS="--sequencer --sequencer-max-safe-lag $SEQUENCER_MAX_SAFE_LAG --sequencer-pk-file $SEQUENCER_PK_FILE"
fi

# Consolidate flags into an array
FLAGS=("--network $NETWORK"
        "--l1-rpc-url $L1_RPC_URL"
        "--l2-rpc-url $L2_RPC_URL"
        "--sync-mode $SYNC_MODE"
        "--l2-engine-url $L2_ENGINE_URL"
        "--jwt-file $JWT_SECRET_PATH"
        "--rpc-port $RPC_PORT"
        "$SYNC_FLAGS" "$DEVNET_FLAGS" "$SEQUENCER_FLAGS" "$@")

echo "Starting sp-magi with the following flags:"
echo "${FLAGS[@]}"

# Start sp-magi with the flags and create wait file
$SP_MAGI_BIN "${FLAGS[@]}" &
PID=$!
echo "PID: $PID"
echo "Creating wait file for docker at $WAITFILE..."
touch $WAITFILE
wait $PID

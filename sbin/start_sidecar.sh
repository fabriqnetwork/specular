#!/bin/bash

set -e

SBIN_DIR=$(dirname "$(readlink -f "$0")")
SBIN_DIR="$(
  cd "$SBIN_DIR"
  pwd
)"
. $SBIN_DIR/utils/utils.sh
ROOT_DIR=$SBIN_DIR/..

require_dotenv "paths" ".paths.env"
require_dotenv "sidecar" ".sidecar.env"

FLAGS=(
  "--l1.endpoint $L1_ENDPOINT"
  "--l2.endpoint $L2_ENDPOINT"
  "--protocol.rollup-cfg-path $ROLLUP_CFG_PATH"
)

if [ "$DISSEMINATOR" = true ]; then
  DISSEMINATOR_PRIV_KEY=$(cat "$DISSEMINATOR_PK_PATH")
  FLAGS+=(
    "--disseminator"
    "--disseminator.private-key $DISSEMINATOR_PRIV_KEY"
    "--disseminator.sub-safety-margin $DISSEMINATOR_SUB_SAFETY_MARGIN"
    "--disseminator.target-batch-size $DISSEMINATOR_TARGET_BATCH_SIZE"
    "--disseminator.max-safe-lag $DISSEMINATOR_MAX_SAFE_LAG"
    "--disseminator.max-safe-lag-delta $DISSEMINATOR_MAX_SAFE_LAG_DELTA"
  )
fi

if [ "$VALIDATOR" = true ]; then
  VALIDATOR_PRIV_KEY=$(cat "$VALIDATOR_PK_PATH")
  FLAGS+=(
    "--validator"
    "--validator.private-key $VALIDATOR_PRIV_KEY"
  )
fi

echo "Executing: $SIDECAR_BIN \${FLAGS[@]}"  # Logging the command to be executed
$SIDECAR_BIN ${FLAGS[@]}  # Executing the command

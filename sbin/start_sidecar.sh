#!/bin/bash

# the local sbin paths are relative to the project root
SBIN=$(dirname "$(readlink -f "$0")")
SBIN="$(
  cd "$SBIN"
  pwd
)"
. $SBIN/utils/utils.sh
ROOT_DIR=$SBIN/..

# Check that the all required dotenv files exists.
reqdotenv "paths" ".paths.env"
reqdotenv "sidecar" ".sidecar.env"

FLAGS=(
  "--l1.endpoint $L1_ENDPOINT"
  "--l2.endpoint $L2_ENDPOINT"
  "--protocol.rollup-cfg-path $ROLLUP_CFG_PATH"
)

# Set disseminator flags.
if [ "$DISSEMINATOR" = true ]; then
  echo "Enabling disseminator."
  DISSEMINATOR_PRIV_KEY=$(cat "$DISSEMINATOR_PK_PATH")
  FLAGS+=(
    "--disseminator"
    "--disseminator.private-key $DISSEMINATOR_PRIV_KEY"
    "--disseminator.sub-safety-margin $DISSEMINATOR_SUB_SAFETY_MARGIN"
    "--disseminator.target-batch-size $DISSEMINATOR_TARGET_BATCH_SIZE"
    "--disseminator.max-safe-lag $DISSEMINATOR_MAX_SAFE_LAG"
    "--disseminator.max-safe-lag-delta $DISSEMINATOR_MAX_SAFE_LAG_DELTA"
  )
  if [ -n "$DISSEMINATOR_INTERVAL" ]; then
    FLAGS+=("--disseminator.interval $DISSEMINATOR_INTERVAL")
  fi
fi
# Set validator flags.
if [ "$VALIDATOR" = true ]; then
  echo "Enabling validator."
  VALIDATOR_PRIV_KEY=$(cat "$VALIDATOR_PK_PATH")
  FLAGS+=(
    "--validator"
    "--validator.private-key $VALIDATOR_PRIV_KEY"
  )
fi

echo "starting sidecar with the following flags:"
echo "${FLAGS[@]}"
$SIDECAR_BIN ${FLAGS[@]}

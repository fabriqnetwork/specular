#!/bin/bash

# the local sbin paths are relative to the project root
SBIN=$(dirname "$(readlink -f "$0")")
SBIN="$(
  cd "$SBIN"
  pwd
)"
ROOT_DIR=$SBIN/..

# Check that the all required dotenv files exists.
PATHS_ENV=".paths.env"
if ! test -f "$PATHS_ENV"; then
  echo "Expected dotenv at $PATHS_ENV (does not exist)."
  exit
fi
echo "Using paths dotenv: $PATHS_ENV"
. $PATHS_ENV
# Use sidecar dotenv (to get l1 endpoint)
SIDECAR_ENV=".sidecar.env"
if ! test -f "$SIDECAR_ENV"; then
    echo "Expected dotenv at $SIDECAR_ENV (does not exist)."
    exit
fi
echo "Using sidecar dotenv: $SIDECAR_ENV"
. $SIDECAR_ENV

###### PID handling ######
trap ctrl_c INT

# Active PIDs
PIDS=()

function cleanup() {
  echo "[run] Cleaning up..."
  for pid in "${PIDS[@]}"; do
    echo "Killing $pid"
    disown $pid
    kill $pid
  done
  # For good measure...
  pgrep geth | xargs kill
  pgrep sidecar | xargs kill
  pgrep magi | xargs kill
  L1_PORT=4012
  if [ -n "$L1_PORT" ]; then
    L1_PORT_PID=$(lsof -i tcp:${L1_PORT} | awk 'NR!=1 {print $2}')
    if [ -n "$L1_WS_PID" ]; then
      echo "Killing proc on $L1_PORT"
      kill $L1_PORT_PID
    fi
  fi
  # Clean up
  $SBIN/clean.sh
}

function ctrl_c() {
  cleanup
}

##########################

WORKSPACE_DIR=./workspace-test
mkdir -p $WORKSPACE_DIR
cd $WORKSPACE_DIR

echo "Cleaning $WORKSPACE_DIR"
$SBIN/clean.sh
# Copy config files to cwd.
echo "Copying local_devnet config files to cwd..."
cp -a $CONFIG_DIR/local_devnet/. .

# Start L1
yes | $SBIN/start_l1.sh -d -s &
L1_PID=$!
PIDS+=($L1_PID)

# Parse url into host:port
L1_HOST_AND_PORT=${L1_ENDPOINT#*://}
# Wait for services
$SBIN/wait-for-it.sh -t 60 $L1_HOST_AND_PORT | sed "s/^/[WAIT] /"
echo "L1 endpoint is available"

# TODO: remove
echo "sleeping for a bit"
sleep 60
echo "done sleeping"

# Start sp-geth
$SBIN/start_sp_geth.sh -c &
SP_GETH_PID=$!
PIDS+=($SP_GETH_PID)

sleep 1

# Start sp-magi
$SBIN/start_sp_magi.sh &
SP_MAGI_PID=$!
PIDS+=($SP_MAGI_PID)

sleep 1

# Start sidecar
$SBIN/start_sidecar.sh &
SIDECAR_PID=$!
PIDS+=($SIDECAR_PID)

cd $CONTRACTS_DIR
echo "Running test $1"
# Run testing script
case $1 in
transactions)
  npx hardhat run scripts/e2e/test_transactions.ts
  RESULT=$?
  ;;
deposit)
  npx hardhat run scripts/e2e/bridge/test_standard_bridge_deposit_eth.ts
  RESULT=$?
  ;;
withdraw)
  npx hardhat run scripts/e2e/bridge/test_standard_bridge_withdraw_eth.ts
  RESULT=$?
  ;;
erc20)
  npx hardhat run scripts/e2e/bridge/test_standard_bridge_erc20.ts
  RESULT=$?
  ;;
*)
  echo "unknown test"
  RESULT=1
  ;;
esac

cleanup
exit $RESULT

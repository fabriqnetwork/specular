#!/bin/bash
set -e

# Check if a test name is provided
if [[ $# -eq 0 ]]; then
  echo "No test name provided"
  exit 1
fi

# Set project root directory
SBIN=$(dirname "$(readlink -f "$0")")
SBIN="$(
  cd "$SBIN"
  pwd
)"
. $SBIN/utils/utils.sh
ROOT_DIR=$SBIN/..

trap ctrl_c INT

# Define cleanup function
function cleanup() {
  echo "Cleaning up processes..."
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

# Handle interrupt signal
function ctrl_c() {
  cleanup
  exit
}

# Check that all required dotenv files exist
reqdotenv "paths" ".paths.env"
reqdotenv "sidecar" ".sidecar.env"

# Set workspace directory
WORKSPACE_DIR=./workspace-test
mkdir -p $WORKSPACE_DIR
cd $WORKSPACE_DIR

# Clean workspace directory
echo "Cleaning $WORKSPACE_DIR"
$SBIN/clean.sh

# Copy config files to cwd
echo "Copying local_devnet config files to cwd..."
cp -a $CONFIG_DIR/e2e_test/. .

# Start L1
yes | $SBIN/generate_secrets.sh -d
yes | $SBIN/start_l1.sh -d -s &

# Parse URL into host:port
L1_HOST_AND_PORT=${L1_ENDPOINT#*://}

# Wait for services
$SBIN/wait-for-it.sh -t 60 $L1_HOST_AND_PORT | sed "s/^/[WAIT] /"
echo "L1 endpoint is available"
until [ -f ".deployed" ]; do
  echo "Waiting for L1 to be fully deployed..."
  sleep 4
done

# Check that deployments dotenv file exists
reqdotenv "deployments" ".deployments.env"

# Start sp-geth
$SBIN/start_sp_geth.sh -c &>proc.out &
sleep 1

# Start sp-magi
$SBIN/start_sp_magi.sh &>proc2.out &
sleep 1

# Start sidecar
$SBIN/start_sidecar.sh &>proc3.out &
sleep 1

# Set contracts directory
cd $CONTRACTS_DIR

# Run testing script based on the provided test name
# Define logging function
function log() {
  echo "$(date) $1"
}

log "Starting the testing script..."

# Run testing script based on the provided test name
log "Running test: $1"
case $1 in
transactions)
  npx hardhat run scripts/e2e/test_transactions.ts
  ;;
deposit)
  npx hardhat run scripts/e2e/bridge/test_standard_bridge_deposit_eth.ts
  ;;
withdraw)
  npx hardhat run scripts/e2e/bridge/test_standard_bridge_withdraw_eth.ts
  ;;
erc20)
  npx hardhat run scripts/e2e/bridge/test_standard_bridge_erc20.ts
  ;;
*)
  log "Unknown test"
  exit 1
  ;;
esac

log "Completed the testing script."

# Clean up and exit with the result of the testing script
cleanup

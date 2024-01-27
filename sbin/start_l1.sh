#!/bin/bash
set -e

# Get the absolute path of the script directory
SBIN=$(dirname "$(readlink -f "$0")")
SBIN="$(
  cd "$SBIN"
  pwd
)"
# Sourcing the utility functions
. $SBIN/utils/utils.sh
ROOT_DIR=$SBIN/..

# Check that all required dotenv files exist
require_dotenv "paths" ".paths.env"
require_dotenv "genesis" ".genesis.env"
require_dotenv "contracts" ".contracts.env"

# Generate waitfile for service init (docker/k8)
# Update the waitfile path if WAIT_DIR is set
WAITFILE="/tmp/.${0##*/}.lock"
if [[ ! -z ${WAIT_DIR+x} ]]; then
  WAITFILE=$WAIT_DIR/.${0##*/}.lock
fi
echo "Using dir $WAIT_DIR for $WAITFILE"

# Remove the waitfile if it exists
if test -f $WAITFILE; then
  rm $WAITFILE
  echo "Removed $WAITFILE"
fi

# Set default values for flags
AUTO_ACCEPT=false
AUTO_APPROVE=""

# Parse arguments
optspec="cdswy"
while getopts "$optspec" optchar; do
  case "${optchar}" in
  c)
    echo "Cleaning deployment before starting l1..."
    $SBIN/clean_deployment.sh
    ;;
  d)
    L1_DEPLOY=true
    ;;
  w)
    WAIT=true
    ;;
  s)
    SILENT=true
    ;;
  y)
    AUTO_ACCEPT=true
    ;;
  *)
    echo "usage: $0 [-c][-d][-s][-y][-h]"
    echo "-c : clean before running"
    echo "-d : deploy contracts"
    echo "-s : silent-mode (no log tailing)"
    echo "-y : auto accept prompts"
    echo "-w : generate docker-compose wait for file"
    exit
    ;;
  esac
done

trap ctrl_c INT

# Set the APPROVE_FLAG if AUTO_ACCEPT is true
if [[ $AUTO_ACCEPT = 'true' ]]; then
  APPROVE_FLAG="-y"
fi

# Parse the L1_HOST and L1_PORT from L1_ENDPOINT
L1_HOST=$(echo $L1_ENDPOINT | awk -F':' '{print substr($2, 3)}')
L1_PORT=$(echo $L1_ENDPOINT | awk -F':' '{print $3}')
echo "Parsed endpoint ($L1_HOST) and port: $L1_PORT from $L1_ENDPOINT"

# Set the log file name
LOG_FILE="l1.log"

# Function to handle interrupt signal
trap ctrl_c INT

# Array to store active PIDs
PIDS=()

# Function to cleanup l1 processes
function cleanup() {
  echo "Cleaning up l1 processes..."
  rm -f $LOG_FILE
  for pid in "${PIDS[@]}"; do
    echo "Killing $pid"
    disown $pid
    kill $pid
  done

  # Remove WAITFILE if WAIT is true
  if [ "$WAIT" = "true" ]; then
    if test -f $WAITFILE; then
      echo "Removing wait file for docker..."
      rm $WAITFILE
    fi
  fi

  # Kill the process on L1_PORT if it exists
  if [ -n "$L1_PORT" ]; then
    L1_WS_PID=$(lsof -i tcp:${L1_PORT} | awk 'NR!=1 {print $2}')
    if [ -n "$L1_WS_PID" ]; then
      echo "Killing proc on $L1_PORT"
      kill $L1_WS_PID
    fi
  fi
}

# Function to handle interrupt signal
function ctrl_c() {
  cleanup
}

# Start L1 network based on L1_STACK
echo "Starting L1..."
if [ "$L1_STACK" = "geth" ]; then
  # Start geth with specified options
  $L1_GETH_BIN \
    --dev \
    --dev.period $L1_PERIOD \
    --verbosity 0 \
    --http \
    --http.api eth,web3,net \
    --http.addr 0.0.0.0 \
    --ws \
    --ws.api eth,net,web3 \
    --ws.addr 0.0.0.0 \
    --ws.port $L1_PORT &>$LOG_FILE &

  # Wait for 1 block
  echo "Waiting for chain progression..."
  sleep $L1_PERIOD

  L1_PID=$

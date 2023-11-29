#!/bin/bash
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

GENESIS_ENV=".genesis.env"
if ! test -f "$GENESIS_ENV"; then
  echo "Expected dotenv at $GENESIS_ENV (does not exist)."
  exit
fi
echo "Using genesis dotenv: $GENESIS_ENV"
. $GENESIS_ENV

if [ "$L1_STACK" = "geth" ]; then
  CONTRACTS_ENV=".contracts.env"
  if ! test -f "$CONTRACTS_ENV"; then
    echo "Expected contracts dotenv at $CONTRACTS_ENV (does not exist)."
    exit
  fi
  echo "Using contracts dotenv: $CONTRACTS_ENV"
  . $CONTRACTS_ENV
fi

# Parse args.
optspec="cdsh"
while getopts "$optspec" optchar; do
  case "${optchar}" in
  c)
    echo "Cleaning deployment before starting l1..."
    $SBIN/clean_deployment.sh
    ;;
  d)
    L1_DEPLOY=true
    ;;
  s)
    SILENT=true
    ;;
  h)
    echo "usage: $0 [-c][-d][-s][-h]"
    echo "-c : clean before running"
    echo "-d : deploy contracts"
    echo "-s : silent-mode (no log tailing)"
    exit
    ;;
  *)
    if [ "$OPTERR" != 1 ] || [ "${optspec:0:1}" = ":" ]; then
      echo "Unknown option: '-${OPTARG}'"
      exit 1
    fi
    ;;
  esac
done

L1_HOST=$(echo $L1_ENDPOINT | awk -F':' '{print substr($2, 3)}')
L1_PORT=$(echo $L1_ENDPOINT | awk -F':' '{print $3}')
echo "Parsed endpoint ($L1_HOST) and port: $L1_PORT from $L1_ENDPOINT"

LOG_FILE="l1.log"

###### PID handling ######
trap ctrl_c INT

# Active PIDs
PIDS=()

function cleanup() {
  echo "Cleaning up l1 processes..."
  rm -f $LOG_FILE
  for pid in "${PIDS[@]}"; do
    echo "Killing $pid"
    disown $pid
    kill $pid
  done
  # For good measure...
  if [ -n "$L1_PORT" ]; then
    L1_WS_PID=$(lsof -i tcp:${L1_PORT} | awk 'NR!=1 {print $2}')
    if [ -n "$L1_WS_PID" ]; then
      echo "Killing proc on $L1_PORT"
      kill $L1_WS_PID
    fi
  fi
}

function ctrl_c() {
  cleanup
}
##########################

# Start L1 network.
echo "Starting L1..."
if [ "$L1_STACK" = "geth" ]; then
  $L1_GETH_BIN \
    --dev \
    --dev.period 4 \
    --http \
    --http.api eth,web3,net \
    --http.addr 0.0.0.0 \
    --ws \
    --ws.api eth,net,web3 \
    --ws.addr 0.0.0.0 \
    --ws.port $L1_PORT &>$LOG_FILE &

  # Wait for 1 block
  echo "Waiting for chain progression..."
  sleep 4

  L1_PID=$!
  echo "L1 PID: $L1_PID"

  echo "Funding addresses..."
  $L1_GETH_BIN attach --exec \
    "eth.sendTransaction({ from: eth.coinbase, to: '"$SEQUENCER_ADDRESS"', value: web3.toWei(10000, 'ether') })" \
    $L1_ENDPOINT
  $L1_GETH_BIN attach --exec \
    "eth.sendTransaction({ from: eth.coinbase, to: '"$VALIDATOR_ADDRESS"', value: web3.toWei(10000, 'ether') })" \
    $L1_ENDPOINT
  $L1_GETH_BIN attach --exec \
    "eth.sendTransaction({ from: eth.coinbase, to: '"$DEPLOYER_ADDRESS"', value: web3.toWei(10000, 'ether') })" \
    $L1_ENDPOINT

  # Wait for 1 block
  echo "Waiting for chain progression..."
  sleep 4
elif [ "$L1_STACK" = "hardhat" ]; then
  echo "Using $CONTRACTS_DIR as HH proj"
  cd $CONTRACTS_DIR && npx hardhat node --no-deploy --hostname $L1_HOST --port $L1_PORT &>$LOG_FILE &
  L1_PID=$!
  PIDS+=($L1_PID)
  echo "L1 PID: $L1_PID"
  sleep 3
else
  echo "invalid value for L1_STACK: $L1_STACK"
  exit 1
fi

# Optionally deploy the contracts
if [ "$L1_DEPLOY" = "true" ]; then
  echo "Deploying contracts..."
  $SBIN/deploy_l1_contracts.sh
fi

# Follow output
if [ ! "$SILENT" = "true" ]; then
    echo "L1 started... (Use ctrl-c to stop)"
    tail -f $LOG_FILE
fi
wait $L1_PID

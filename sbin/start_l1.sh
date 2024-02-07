#!/bin/bash
set -e
SBIN=$(dirname "$(readlink -f "$0")")
SBIN="$(
  cd "$SBIN"
  pwd
)"
. $SBIN/utils/utils.sh
ROOT_DIR=$SBIN/..

WAITFILE="/tmp/.${0##*/}.lock"

if [[ ! -z ${WAIT_DIR+x} ]]; then
  WAITFILE=$WAIT_DIR/.${0##*/}.lock
fi

if test -f $WAITFILE; then
  rm $WAITFILE
  echo "Removed $WAITFILE"
fi

WORKSPACE_DIR=$HOME/.spc/workspaces/active_workspace

PATHS_ENV=$WORKSPACE_DIR/.paths.env
GENESIS_ENV=$WORKSPACE_DIR/.genesis.env
CONTRACTS_ENV=$WORKSPACE_DIR/.contracts.env

reqdotenv "paths" $PATHS_ENV
reqdotenv "genesis" $GENESIS_ENV
reqdotenv "contracts" $CONTRACTS_ENV

# Generate waitfile for service init (docker/k8)
WAITFILE="/tmp/.${0##*/}.lock"

if [[ ! -z ${WAIT_DIR+x} ]]; then
  WAITFILE=$WAIT_DIR/.${0##*/}.lock
fi

echo "Using dir $WAIT_DIR for $WAITFILE"

AUTO_ACCEPT=false
AUTO_APPROVE=""
# Parse args.
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
if [[ $AUTO_ACCEPT = 'true' ]]; then
  APPROVE_FLAG="-y"
fi

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

  # Remove WAITFILE
  if [ "$WAIT" = "true" ]; then
    if test -f $WAITFILE; then
      echo "Removing wait file for docker..."
      rm $WAITFILE
    fi
  fi

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
L1_PERIOD=$L1_PERIOD spc up l1geth &

L1_PID=$!
echo "L1 PID: $L1_PID"

# Wait for 1 block
echo "Waiting for chain progression..."
sleep $L1_PERIOD

echo "Funding addresses..."
# Add addresses from .contracts.env
addresses_to_fund=($SEQUENCER_ADDRESS $VALIDATOR_ADDRESS $DEPLOYER_ADDRESS)
# Add addresses from $GENESIS_CFG_PATH
addresses_to_fund+=($(python3 -c "import json; print(' '.join(json.load(open('$GENESIS_CFG_PATH'))['alloc']))"))
# TODO: consider using cast (more general)
for address in "${addresses_to_fund[@]}"; do
  mycall="eth.sendTransaction({ from: eth.coinbase, to: '"$address"', value: web3.toWei(10000, 'ether') })"
  $L1_GETH_BIN attach --exec "$mycall" $L1_ENDPOINT
done
# Wait for 1 block
echo "Waiting for chain progression..."
sleep $L1_PERIOD

# Optionally deploy the contracts
if [ "$L1_DEPLOY" = "true" ]; then
  echo "Deploying contracts..."
  $SBIN/deploy_l1_contracts.sh $APPROVE_FLAG
fi

# Follow output
if [ ! "$SILENT" = "true" ]; then
  if [ "$WAIT" = "true" ]; then
    echo "Creating wait file for docker at $WAITFILE..."
    touch $WAITFILE
  fi
  echo "L1 started... (Use ctrl-c to stop)"
  tail -f $LOG_FILE
fi
wait $L1_PID

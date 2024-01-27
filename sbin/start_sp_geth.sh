#!/bin/bash
set -e

# Set SBIN to the absolute path of the script's directory
SBIN=$(dirname "$(readlink -f "$0")")
SBIN="$(
  cd "$SBIN"
  pwd
)"
# Source the utility functions
. $SBIN/utils/utils.sh
ROOT_DIR=$SBIN/..

# Check that the required dotenv files exist
require_dotenv "paths" ".paths.env"
require_dotenv "sp_geth" ".sp_geth.env"

# Generate waitfile for service init (docker/k8)
WAITFILE="/tmp/.${0##*/}.lock"

# If WAIT_DIR is set, use it as the location for the wait file
if [[ ! -z ${WAIT_DIR+x} ]]; then
  WAITFILE=$WAIT_DIR/.${0##*/}.lock
fi

# Parse options
optspec="chw"
while getopts "$optspec" optchar; do
  case "${optchar}" in
  w)
    WAIT=true
    ;;
  c)
    # Remove debugging statement
    $SBIN/clean_sp_geth.sh
    ;;
  h)
    # Print usage information and exit
    echo "usage: $0 [-c][-h][-w]"
    echo "-c : clean before running"
    echo "-w : generate docker-compose wait file"
    exit
    ;;
  *)
    # Handle unknown options
    if [ "$OPTERR" != 1 ] || [ "${optspec:0:1}" = ":" ]; then
      echo "Unknown option: '-${OPTARG}'"
      exit 1
    fi
    ;;
  esac
done

# Remove docker wait file if WAIT is set to true
if [ "$WAIT" = "true" ]; then
  if test -f $WAITFILE; then
    echo "Removing wait file for docker..."
    rm $WAITFILE
  fi
fi

# Initialize sp-geth with genesis json if necessary
if [ ! -d $DATA_DIR ]; then
  echo "Initializing sp-geth with genesis json at $GENESIS_PATH"
  if [ ! -f $GENESIS_PATH ]; then
    echo "Missing genesis json at $GENESIS_PATH"
    exit 1
  fi
  $SP_GETH_BIN --datadir $DATA_DIR --networkid $NETWORK_ID init $GENESIS_PATH
fi

# Start sp-geth.
FLAGS="
    --datadir $DATA_DIR \
    --networkid $NETWORK_ID \
    --http \
    --http.addr $ADDRESS \
    --http.port $HTTP_PORT \
    --http.api engine,personal,eth,net,web3,txpool,miner,debug \
    --http.corsdomain=* \
    --http.vhosts=* \
    --ws \
    --ws.addr $ADDRESS \
    --ws.port $WS_PORT \
    --ws.api engine,personal,eth,net,web3,txpool,miner,debug \
    --ws.origins=* \
    --authrpc.vhosts=* \
    --authrpc.addr $ADDRESS \
    --authrpc.port $AUTH_PORT \
    --authrpc.jwtsecret $JWT_SECRET_PATH \
    --miner.recommit 0 \
    --nodiscover \
    --maxpeers 0 \
    --syncmode full
"

echo "Starting sp-geth with the following aruments:"
# Start sp-geth with specified flags
echo "Starting sp-geth with the following arguments:"
echo $FLAGS
$SP_GETH_BIN $FLAGS &

PID=$!
echo "PID: $PID"

# Wait for sp-geth to start
sleep 15

# Create wait file for docker if WAIT is set to true
if [ "$WAIT" = "true" ]; then
  echo "Creating wait file for docker at $WAITFILE..."
  touch $WAITFILE
fi

# Wait for sp-geth to finish
wait $PID

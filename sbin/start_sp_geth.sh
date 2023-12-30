#!/bin/bash

# currently the local sbin paths are relative to the project root
SBIN=$(dirname "$(readlink -f "$0")")
SBIN="$(
  cd "$SBIN"
  pwd
)"
. $SBIN/utils/utils.sh
ROOT_DIR=$SBIN/..
SP_GETH_WAIT=/tmp/.sp_geth_started.lock
# Check that the all required dotenv files exists.
reqdotenv "paths" ".paths.env"
reqdotenv "sp_geth" ".sp_geth.env"

# Parse args.
optspec="ch"
while getopts "$optspec" optchar; do
  case "${optchar}" in
  l)
    echo "Creating wait file for docker"
    touch $SP_GETH_WAIT
    ;;
  c)
    echo "Cleaning..."
    $SBIN/clean_sp_geth.sh
    ;;
  h)
    echo "usage: $0 [-c][-h]"
    echo "-c : clean before running"
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

if [ ! -d $DATA_DIR ]; then
  echo "Initializing sp-geth with genesis json at $GENESIS_PATH"
  if [ ! -f $GENESIS_PATH ]; then
    echo "Missing genesis json at $GENESIS_PATH"
    exit
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
echo $FLAGS

# Remove SP_GETH_WAITFILE
if [ "$SP_GETH_WAIT" = "true" ]; then
  echo "Adding wait file for docker..."
    rm $SP_GETH_WAITFILE
fi
$SP_GETH_BIN $FLAGS

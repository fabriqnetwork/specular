#!/bin/bash

# currently the local sbin paths are relative to the project root
SBIN=$(dirname "$(readlink -f "$0")")
SBIN="$(
  cd "$SBIN"
  pwd
)"
. $SBIN/utils/utils.sh
ROOT_DIR=$SBIN/..


# Check that the all required dotenv files exists.
reqdotenv "paths" ".paths.env"
reqdotenv "sp_geth" ".sp_geth.env"

# Generate waitfile for service init (docker/k8)
WAITFILE="/tmp/.${0##*/}.lock"

if [[ ! -z ${WAIT_DIR+x} ]]; then
  WAITFILE=$WAIT_DIR/.${0##*/}.lock
fi

# Parse args.
optspec="chw"
while getopts "$optspec" optchar; do
  case "${optchar}" in
  w)
    echo "Creating wait file for docker"
    touch $WAITFILE
    ;;
  c)
    echo "Cleaning..."
    $SBIN/clean_sp_geth.sh
    ;;
  h)
    echo "usage: $0 [-c][-h][-w]"
    echo "-c : clean before running"
    echo "-w : generate docker-compose wait for file"
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

# Remove WAITFILEFILE
if [ "$WAITFILE" = "true" ]; then
  echo "Removing wait file for docker..."
  rm $WAITFILEFILE
fi
$SP_GETH_BIN $FLAGS

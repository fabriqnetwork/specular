#!/bin/bash

# currently the local sbin paths are relative to the project root
SBIN=$(dirname "$(readlink -f "$0")")
SBIN="`cd "$SBIN"; pwd`"
ROOT_DIR=$SBIN/..

# Check that the all required dotenv files exists.
PATHS_ENV=".paths.env"
if ! test -f "$PATHS_ENV"; then
    echo "Expected dotenv at $PATHS_ENV (does not exist)."
    exit
fi
echo "Using paths dotenv: $PATHS_ENV"
. $PATHS_ENV

SP_GETH_ENV=".sp_geth.env"
if ! test -f "$SP_GETH_ENV"; then
    echo "Expected dotenv at $SP_GETH_ENV (does not exist)."
    exit
fi
echo "Using dotenv: $SP_GETH_ENV"
. $SP_GETH_ENV

# Parse args.
optspec="ch"
while getopts "$optspec" optchar; do
    case "${optchar}" in
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
    echo "Initializing sp-geth..."
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
$SP_GETH_BIN $FLAGS

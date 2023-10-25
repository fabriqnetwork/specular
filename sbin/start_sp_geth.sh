#!/bin/bash
# TODO: rename this file to start_sp_geth.sh to disambiguate
if [ -z $SP_GETH ]; then
    # If no binary specified, assume repo directory structure.
    SBIN=`dirname $0`
    SBIN="`cd "$SBIN"; pwd`"
    . $SBIN/configure.sh
    SP_GETH=$GETH_BIN
fi
echo "Using bin: $SP_GETH"

# Check that the dotenv exists.
ENV=".sp_geth.env"
if ! test -f $ENV; then
    echo "Expected dotenv at $ENV (does not exist)."
    exit
fi
echo "Using dotenv: $ENV"
. $ENV

if [ ! -d $DATA_DIR ]; then
    echo "Initializing sp-geth..."
    $SP_GETH --datadir $DATA_DIR --networkid $NETWORK_ID init $GENESIS_PATH
fi

# Start sp-geth.
args=(
    --datadir $DATA_DIR
    --networkid $NETWORK_ID
    --http
    --http.addr "0.0.0.0"
    --http.port $HTTP_PORT
    --http.api "engine,personal,eth,net,web3,txpool,miner,debug"
    --http.corsdomain "*"
    --http.vhosts "*"
    --ws
    --ws.addr "0.0.0.0"
    --ws.port $WS_PORT
    --ws.api "engine,personal,eth,net,web3,txpool,miner,debug"
    --ws.origins "*" \
    --authrpc.vhosts "*" \
    --authrpc.addr 0.0.0.0 \
    --authrpc.port $AUTH_PORT \
    --authrpc.jwtsecret $JWT_SECRET_PATH \
    --syncmode=full \
    --miner.recommit 0
)
echo "Starting sp-geth..."
$SP_GETH "${args[@]}"

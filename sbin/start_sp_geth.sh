#!/bin/bash
SBIN=`dirname $0`
SBIN="`cd "$SBIN"; pwd`"
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

# Check that the dotenv exists.
ENV=".sp_geth.env"
if ! test -f $ENV; then
    echo "Expected dotenv at $ENV (does not exist)."
    exit
fi
echo "Using dotenv: $ENV"
. $ENV

if [ -z $SP_GETH_BIN ]; then
    # If no binary specified, assume repo directory structure.
    . $SBIN/configure.sh
fi
echo "Using bin: $SP_GETH_BIN"

if [ ! -d $DATA_DIR ]; then
    echo "Initializing sp-geth..."
    $SP_GETH_BIN --datadir $DATA_DIR --networkid $NETWORK_ID init $GENESIS_PATH
fi

# Start sp-geth.
args="
    --datadir $DATA_DIR \
    --networkid $NETWORK_ID \
    --http \
    --http.addr $ADDR \
    --http.port $HTTP_PORT \
    --http.api 'engine,personal,eth,net,web3,txpool,miner,debug' \
    --http.corsdomain=* \
    --http.vhosts=* \
    --ws \
    --ws.addr $ADDR \
    --ws.port $WS_PORT \
    --ws.api 'engine,personal,eth,net,web3,txpool,miner,debug' \
    --ws.origins=* \
    --authrpc.vhosts=* \
    --authrpc.addr $ADDR \
    --authrpc.port $AUTH_PORT \
    --authrpc.jwtsecret $JWT_SECRET_PATH \
    --miner.recommit 0 \
"

echo "Starting sp-geth with the following aruments:"
echo $args

$SP_GETH_BIN $args

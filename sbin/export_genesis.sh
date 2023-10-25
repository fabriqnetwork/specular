#!/bin/sh

# Check that the dotenv exists, or GENESIS_PATH is set.
ENV=".genesis.env"
if ! test -f $ENV && [ -z ${GENESIS_PATH+x} ]; then
    echo "Expected GENESIS_PATH (not set) OR dotenv at $ENV (does not exist)."
    exit 1
fi
. $ENV

# Note: do not add any printing to stdout in the positive case.
if [ -z $SP_GETH ]; then
    # If no binary specified, assume repo directory structure.
    SBIN=`dirname $0`
    SBIN="`cd "$SBIN"; pwd`"
    . $SBIN/configure.sh
    SP_GETH=$GETH_BIN
fi

# Export l2 genesis hash for $GENESIS_PATH
DATA_DIR=tmp_data/
HTTP_ADDRESS="0.0.0.0"
HTTP_PORT=1234
# Initialize sp-geth
$SP_GETH init --datadir $DATA_DIR $GENESIS_PATH
# Start sp-geth
$SP_GETH --datadir $DATA_DIR --http --http.addr $HTTP_ADDRESS --http.port 1234 &
SP_GETH_PID=$!
sleep 1
# Export l2 genesis hash for $GENESIS_PATH
RESULT=`$SP_GETH attach --exec "eth.getBlock(0).hash" "http://$HTTP_ADDRESS:$HTTP_PORT"`
kill $SP_GETH_PID
rm -r $DATA_DIR

echo "{\"hash\": $RESULT}"

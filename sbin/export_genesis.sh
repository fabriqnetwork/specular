#!/bin/bash
SBIN=$(dirname "$(readlink -f "$0")")
ROOT_DIR=$SBIN/..

cd $ROOT_DIR/workspace

# Check that the all required dotenv files exists.
PATHS_ENV=".paths.env"
if ! test -f "$PATHS_ENV"; then
  echo "Expected dotenv at $PATHS_ENV (does not exist)."
  exit
fi
. $PATHS_ENV

# Check that the dotenv exists, or GENESIS_PATH is set.
ENV=".genesis.env"
if ! test -f $ENV && [ -z ${GENESIS_PATH+x} ]; then
  echo "Expected GENESIS_PATH (not set) OR dotenv at $ENV (does not exist)."
  exit 1
fi
. $ENV

# Export l2 genesis hash for $GENESIS_PATH
DATA_DIR=tmp_data/
HTTP_ADDRESS="0.0.0.0"
HTTP_PORT=1234
# Initialize sp-geth
$SP_GETH_BIN init --datadir $DATA_DIR $GENESIS_PATH
# Start sp-geth
$SP_GETH_BIN --datadir $DATA_DIR --http --http.addr $HTTP_ADDRESS --http.port 1234 &
SP_GETH_PID=$!
sleep 1
# Export l2 genesis hash for $GENESIS_PATH
RESULT=$($SP_GETH_BIN attach --exec "eth.getBlock(0).hash" "http://$HTTP_ADDRESS:$HTTP_PORT")
kill $SP_GETH_PID
rm -r $DATA_DIR

echo "{\"hash\": $RESULT}"

#!/bin/bash

# Configure variables
SBIN_DIR=`dirname $0`
SBIN_DIR="`cd "$SBIN_DIR"; pwd`"
set -o allexport
source $SBIN_DIR/configure.sh
set +o allexport

# Spin up L1 node
cd $CONTRACTS_DIR
ganache --chain.chainId 31337 -b 5 -m "test test test test test test test test test test test junk" 2>&1 &
GANACHE_PID=$!
npx hardhat deploy --network localhost

# Spin up L2 node
cd $PROJECT_DATA_DIR
$SBIN_DIR/sequencer.sh > $PROJECT_LOG_DIR/l2.log 2>&1 &
L2GETH_PID=$!

# Wait for nodes
$SBIN_DIR/wait-for-it.sh -t 60 $HOST:$L1_WS_PORT
$SBIN_DIR/wait-for-it.sh -t 60 $HOST:$L2_HTTP_PORT

# Run testing script
cd $CONTRACTS_DIR
npx hardhat deploy --network specularLocalDev
sleep 10

case $1 in

  general)
    npx ts-node scripts/testing.ts
    RESULT=$?
    ;;

  deposit)
    npx hardhat run scripts/bridge/standard_bridge_deposit_eth.ts
    RESULT=$?
    ;;

  withdraw)
    npx hardhat run scripts/bridge/standard_bridge_withdraw_eth.ts
    RESULT=$?
    ;;

  erc20)
    npx hardhat run scripts/bridge/standard_bridge_erc20.ts
    RESULT=$?
    ;;

  *)
    echo "unknown test"
    ;;
esac


# Kill nodes
disown $L2GETH_PID
disown $GANACHE_PID
kill $L2GETH_PID
kill $GANACHE_PID

# Clean up
$SBIN_DIR/clean.sh

exit $RESULT

#!/bin/bash

# the local sbin paths are relative to the project root
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

###### PID handling ######
trap ctrl_c INT

# Active PIDs
PIDS=()

function cleanup() {
    echo "Cleaning up..."
    for pid in "${PIDS[@]}"; do
	echo "Killing $pid"
	disown $pid
	kill $pid
    done
    # Clean up
    $SBIN/clean.sh
}

function ctrl_c() {
    cleanup
}

##########################

$SBIN/clean.sh
# Copy config files to cwd.
echo "Copying local_devnet config files to cwd..."
cp -a $CONFIG/deployments/local_devnet/. .

# Use sidecar .env (to get l1 endpoint)
ENV=".sidecar.env"
echo "Using dotenv: $ENV"
. $ENV
# Parse url into host:port for wait-for-it.sh
L1_HOST_AND_PORT=${L1_ENDPOINT#*://}

# TODO: improve logs accross these scripts
# Start L1
$SBIN/start_l1.sh &
L1_PID=$!
PIDS+=$L1_PID

# Start sidecar
$SBIN/start_sidecar.sh &
SIDECAR_PID=$!
PIDS+=$SIDECAR_PID
# Start sp-geth
$SBIN/start_sp_geth.sh &
SP_GETH_PID=$!
PIDS+=$SP_GETH_PID

# Wait for services
$SBIN/wait-for-it.sh -t 60 $L1_HOST_AND_PORT | sed "s/^/[WAIT] /"
$SBIN/wait-for-it.sh -t 60 $L1_HOST_AND_PORT | sed "s/^/[WAIT] /"

cd $CONTRACTS_DIR
npx hardhat deploy --network specularLocalDev | sed "s/^/[L2 deploy] /"

# Run testing script
case $1 in
  transactions)
    npx hardhat run scripts/e2e/test_transactions.ts
    RESULT=$?
    ;;
  deposit)
    npx hardhat run scripts/e2e/bridge/test_standard_bridge_deposit_eth.ts
    RESULT=$?
    ;;
  withdraw)
    npx hardhat run scripts/e2e/bridge/test_standard_bridge_withdraw_eth.ts
    RESULT=$?
    ;;
  erc20)
    npx hardhat run scripts/e2e/bridge/test_standard_bridge_erc20.ts
    RESULT=$?
    ;;
  *)
    echo "unknown test"
    RESULT=1
    ;;
esac

cleanup
exit $RESULT

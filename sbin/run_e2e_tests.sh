#!/bin/bash
SBIN=$(dirname "$(readlink -f "$0")")
SBIN="`cd "$SBIN"; pwd`"
ROOT="`cd $SBIN/../; pwd`"
CONFIG="$ROOT/config"
. $SBIN/configure.sh

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
$SBIN/start_l1.sh &
# TODO: this is not actually working right now
$SBIN/start_sidecar.sh &
SIDECAR_PID=$!
echo "sidecar PID=$SIDECAR_PID"
$SBIN/start_geth.sh &
SP_GETH_PID=$!
echo "sp-geth PID=$SP_GETH_PID"

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


# Kill nodes
disown $SP_GETH_PID
disown $SIDECAR_PID
kill $SP_GETH_PID
kill $SIDECAR_PID

# Clean up
$SBIN/clean.sh

exit $RESULT

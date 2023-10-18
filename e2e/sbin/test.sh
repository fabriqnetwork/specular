#!/bin/bash
E2E_SBIN_DIR=$(dirname "$(readlink -f "$0")")
E2E_SBIN_DIR="`cd "$E2E_SBIN_DIR"; pwd`"

cd $E2E_SBIN_DIR/../../sbin
. ./configure.sh

# TODO: improve logs accross these scripts

$SBIN/clean.sh
$SBIN/start_l1.sh 2>&1 | sed "s/^/[L1] /"

$SBIN/init_ripcord.sh 2>&1 | sed "s/^/[L2] /"
$SBIN/start_ripcord.sh 2>&1 | sed "s/^/[L2] /" &

sleep 5

# Wait for nodes
$E2E_SBIN_DIR/wait-for-it.sh -t 60 $HOST:$L1_WS_PORT | sed "s/^/[WAIT] /"
$E2E_SBIN_DIR/wait-for-it.sh -t 60 $HOST:$L2_HTTP_PORT | sed "s/^/[WAIT] /"

sleep 5

cd $CONTRACTS_DIR
npx hardhat deploy --network specularLocalDev | sed "s/^/[L2 deploy] /"
# docker logs geth_container --follow

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
disown $GETH_PID
disown $SIDECAR_PID
kill $GETH_PID
kill $SIDECAR_PID

# Clean up
$SBIN_DIR/clean.sh

exit $RESULT

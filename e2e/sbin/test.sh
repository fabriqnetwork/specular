#!/bin/bash

# Configure variables
SBIN_DIR=`dirname $0`
SBIN_DIR="`cd "$SBIN_DIR"; pwd`"
set -o allexport
source $SBIN_DIR/configure.sh
set +o allexport

# Spin up L1 node
cd $CONTRACTS_DIR
docker run -d \
  --name geth_container \
  -v ./../docker:/root \
  -p 8545:8545 \
  ethereum/client-go \
  --dev --dev.period 1 \
  --verbosity 3 \
  --http --http.api eth,web3,net --http.addr 0.0.0.0 \
  --ws --ws.api eth,net,web3 --ws.addr 0.0.0.0 --ws.port 8545 2>&1 | sed "s/^/[L1] /"

sleep 3

GETH_DOCKER_URL="ws://172.17.0.1:8545"

docker exec geth_container geth attach --exec \
  "eth.sendTransaction({ from: eth.coinbase, to: '"$SEQUENCER_ADDR"', value: web3.toWei(10000, 'ether') })" \
  $GETH_DOCKER_URL | sed "s/^/[L1] /"

docker exec geth_container geth attach --exec \
  "eth.sendTransaction({ from: eth.coinbase, to: '"$VALIDATOR_ADDR"', value: web3.toWei(10000, 'ether') })" \
  $GETH_DOCKER_URL | sed "s/^/[L1] /"

docker exec geth_container geth attach --exec \
  "eth.sendTransaction({ from: eth.coinbase, to: '"$DEPLOYER_ADDR"', value: web3.toWei(10000, 'ether') })" \
  $GETH_DOCKER_URL | sed "s/^/[L1] /"

sleep 10

npx hardhat deploy --network localhost | sed "s/^/[L1] /"

# Spin up L2 node
cd $PROJECT_DATA_DIR
$SBIN_DIR/sequencer.sh 2>&1 | sed "s/^/[L2] /" &
L2GETH_PID=$!

# Wait for nodes
$SBIN_DIR/wait-for-it.sh -t 60 $HOST:$L1_WS_PORT
$SBIN_DIR/wait-for-it.sh -t 60 $HOST:$L2_HTTP_PORT

# Run testing script
cd $CONTRACTS_DIR
npx hardhat deploy --network specularLocalDev | sed "s/^/[L2] /"

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
disown $L2GETH_PID
kill $L2GETH_PID

# Clean up
$SBIN_DIR/clean.sh

exit $RESULT

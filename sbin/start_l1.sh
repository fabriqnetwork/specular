#!/bin/bash
SBIN_DIR=`dirname $0`
SBIN_DIR="`cd "$SBIN_DIR"; pwd`"
source $SBIN_DIR/configure.sh

cd $ROOT_DIR

docker remove --force geth_container

docker run -d \
  --name geth_container \
  -v ./docker:/root \
  -p 8545:8545 \
  ethereum/client-go \
  --dev --dev.period 1 \
  --verbosity 3 \
  --http --http.api eth,web3,net --http.addr 0.0.0.0 \
  --ws --ws.api eth,net,web3 --ws.addr 0.0.0.0 --ws.port 8545

sleep 3

GETH_DOCKER_URL="ws://172.17.0.1:8545"

docker exec geth_container geth attach --exec \
  "eth.sendTransaction({ from: eth.coinbase, to: '"$SEQUENCER_ADDR"', value: web3.toWei(10000, 'ether') })" \
  $GETH_DOCKER_URL

docker exec geth_container geth attach --exec \
  "eth.sendTransaction({ from: eth.coinbase, to: '"$VALIDATOR_ADDR"', value: web3.toWei(10000, 'ether') })" \
  $GETH_DOCKER_URL

docker exec geth_container geth attach --exec \
  "eth.sendTransaction({ from: eth.coinbase, to: '"$DEPLOYER_ADDR"', value: web3.toWei(10000, 'ether') })" \
  $GETH_DOCKER_URL

docker exec geth_container geth attach --exec \
  "eth.sendTransaction({ from: eth.coinbase, to: '"$BRIDGER_ADDR"', value: web3.toWei(10000, 'ether') })" \
  $GETH_DOCKER_URL | sed "s/^/[L1] /"

cd $CONTRACTS_DIR
npx hardhat deploy --network localhost

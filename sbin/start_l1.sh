#!/bin/bash
if [ -z $CONTRACTS_DIR ]; then
    SBIN=`dirname $0`
    SBIN="`cd "$SBIN"; pwd`"
    . $SBIN/configure.sh
fi
echo "Using $CONTRACTS_DIR as deployment source"

# Check that the dotenv exists.
ENV=".sidecar.env"
if ! test -f $ENV; then
    echo "Expected dotenv at $ENV (does not exist)."
    exit
fi
echo "Using dotenv: $ENV"
. $ENV


docker rm --force geth_container

# Start L1 network.
docker run -d \
  --name geth_container \
  -v ./docker:/root \
  -p $L1_WS_PORT:$L1_WS_PORT \
  ethereum/client-go \
  --dev \
  --dev.period 1 \
  --verbosity 3 \
  --http \
  --http.api eth,web3,net \
  --http.addr 0.0.0.0 \
  --ws \
  --ws.api eth,net,web3 \
  --ws.addr 0.0.0.0 \
  --ws.port $L1_WS_PORT

sleep 3

# Fund sidecar addresses.
docker exec geth_container geth attach --exec \
  "eth.sendTransaction({ from: eth.coinbase, to: '"$SEQUENCER_ADDR"', value: web3.toWei(10000, 'ether') })" \
  $L1_ENDPOINT

docker exec geth_container geth attach --exec \
  "eth.sendTransaction({ from: eth.coinbase, to: '"$VALIDATOR_ADDR"', value: web3.toWei(10000, 'ether') })" \
  $L1_ENDPOINT

docker exec geth_container geth attach --exec \
  "eth.sendTransaction({ from: eth.coinbase, to: '"$DEPLOYER_ADDR"', value: web3.toWei(10000, 'ether') })" \
  $L1_ENDPOINT


# Deploy contracts
cd $CONTRACTS_DIR
npx hardhat deploy --network localhost

docker logs geth_container --follow

#!/bin/bash
if [ -z $CONTRACTS_DIR ]; then
    SBIN=`dirname $0`
    SBIN="`cd "$SBIN"; pwd`"
    . $SBIN/configure.sh
fi
echo "Using $CONTRACTS_DIR as deployment source"

# Check that the dotenv exists.
ENV=".genesis.env"
if ! test -f $ENV; then
    echo "Expected dotenv at $ENV (does not exist)."
    exit
fi
echo "Using dotenv: $ENV"
. $ENV

echo "Force-removing l1_geth container if it exists..."
docker rm --force l1_geth

# Start L1 network.
echo "Starting L1..."
L1_WS_PORT=`echo $L1_ENDPOINT | awk -F':' '{print $3}'`
docker run -d \
  --name l1_geth \
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

echo "Funding addresses..."
docker exec l1_geth geth attach --exec \
  "eth.sendTransaction({ from: eth.coinbase, to: '"$SEQUENCER_ADDRESS"', value: web3.toWei(10000, 'ether') })" \
  $L1_ENDPOINT
docker exec l1_geth geth attach --exec \
  "eth.sendTransaction({ from: eth.coinbase, to: '"$VALIDATOR_ADDRESS"', value: web3.toWei(10000, 'ether') })" \
  $L1_ENDPOINT
docker exec l1_geth geth attach --exec \
  "eth.sendTransaction({ from: eth.coinbase, to: '"$DEPLOYER_ADDRESS"', value: web3.toWei(10000, 'ether') })" \
  $L1_ENDPOINT


# Deploy contracts
echo "Deploying l1 contracts..."
relpath () {
    echo `python3 -c "import os.path; print(os.path.relpath('$1', '$2'))"`
}
GENESIS_PATH=`relpath $GENESIS_PATH $CONTRACTS_DIR` # use relative path
cd $CONTRACTS_DIR
GENESIS_PATH=$GENESIS_PATH npx hardhat deploy --network localhost

docker logs l1_geth --follow

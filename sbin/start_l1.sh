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

echo "Removing localhost deployment artifacts..."
rm -rf $CONTRACTS_DIR/deployments/localhost

CWD=`pwd`
L1_WS_PORT=`echo $L1_ENDPOINT | awk -F':' '{print $3}'`
echo "Parsed port: $L1_WS_PORT from $L1_ENDPOINT"

###### PID handling ######
trap ctrl_c INT

# Active PIDs
PIDS=()

function cleanup() {
    echo "Cleaning up..."
    for pid in "${PIDS[@]}"; do
	echo "Killing $pid"
	kill $pid
    done
}

function ctrl_c() {
    cleanup
}
##########################

# Start L1 network.
echo "Starting L1..."
if [ $L1_STACK = "geth" ]; then
    echo "Force-removing l1_geth container if it exists..."
    docker rm --force l1_geth
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
elif [ $L1_STACK = "hardhat" ]; then
    cd $CONTRACTS_DIR && npx hardhat node --no-deploy --port $L1_WS_PORT &
    L1_PID=$!
    PIDS+=$L1_PID
    echo "L1 PID: $L1_PID"
    sleep 3
else
    echo "invalid value for L1_STACK: $L1_STACK"
    exit 1
fi

# Deploy contracts
echo "Deploying l1 contracts..."
cd $CONTRACTS_DIR && npx hardhat deploy --network localhost

echo "Generating rollup config..."
cd $CWD && $SBIN/create_rollup_config.sh

# Follow output
if [ $L1_STACK = "geth" ]; then
    docker logs l1_geth --follow
elif [ $L1_STACK = "hardhat" ]; then
    tail -f $PID
fi

#!/bin/bash
SBIN=`dirname $0`
SBIN="`cd "$SBIN"; pwd`"
if [ -z $CONTRACTS_DIR ] | [ -z $L1_GETH_BIN ]; then
    . $SBIN/configure.sh
fi
# Parse args.
optspec="cdh"
while getopts "$optspec" optchar; do
    case "${optchar}" in
        c)
	    echo "Cleaning..."
	    $SBIN/clean_deployment.sh
	    ;;
	d)
	    L1_DEPLOY=true
	    ;;
        h)
            echo "usage: $0 [-c][-d][-h]"
	    echo "-c : clean before running"
	    echo "-d : deploy contracts"
            exit
            ;;
        *)
            if [ "$OPTERR" != 1 ] || [ "${optspec:0:1}" = ":" ]; then
                echo "Unknown option: '-${OPTARG}'"
		exit 1
            fi
            ;;
    esac
done

# Check that the dotenv exists.
GENESIS_ENV=".genesis.env"
if ! test -f $GENESIS_ENV; then
    echo "Expected dotenv at $GENESIS_ENV (does not exist)."
    exit
fi
echo "Using genesis dotenv: $GENESIS_ENV"
. $GENESIS_ENV
CONTRACTS_ENV=".contracts.env"
if ! test -f $CONTRACTS_ENV; then
    echo "Expected dotenv at $CONTRACTS_ENV (does not exist)."
    exit
fi
echo "Using contracts dotenv: $CONTRACTS_ENV"
. $CONTRACTS_ENV

L1_HOST=`echo $L1_ENDPOINT | awk -F':' '{print substr($2, 3)}'`
L1_PORT=`echo $L1_ENDPOINT | awk -F':' '{print $3}'`
echo "Parsed endpoint ($L1_HOST) and port: $L1_PORT from $L1_ENDPOINT"

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
    # For good measure...
    if [ -n "$L1_PORT" ]; then
        L1_WS_PID=`lsof -i tcp:${L1_PORT} | awk 'NR!=1 {print $2}'`
        if [ -n "$L1_WS_PID" ]; then
	        echo "Killing proc on $L1_PORT"
            kill $L1_WS_PID
        fi
    fi
}

function ctrl_c() {
    cleanup
}
##########################

# Start L1 network.
echo "Starting L1..."
if [ "$L1_STACK" = "geth" ]; then
    $L1_GETH_BIN \
      --dev \
      --verbosity 0 \
      --http \
      --http.api eth,web3,net \
      --http.addr 0.0.0.0 \
      --ws \
      --ws.api eth,net,web3 \
      --ws.addr 0.0.0.0 \
      --ws.port $L1_PORT &

    L1_PID=$!
    echo "L1 PID: $L1_PID"

    sleep 3

    echo "Funding addresses..."
    $L1_GETH_BIN attach --exec \
      "eth.sendTransaction({ from: eth.coinbase, to: '"$SEQUENCER_ADDRESS"', value: web3.toWei(10000, 'ether') })" \
      $L1_ENDPOINT
    $L1_GETH_BIN attach --exec \
      "eth.sendTransaction({ from: eth.coinbase, to: '"$VALIDATOR_ADDRESS"', value: web3.toWei(10000, 'ether') })" \
      $L1_ENDPOINT
    $L1_GETH_BIN attach --exec \
      "eth.sendTransaction({ from: eth.coinbase, to: '"$DEPLOYER_ADDRESS"', value: web3.toWei(10000, 'ether') })" \
      $L1_ENDPOINT
elif [ "$L1_STACK" = "hardhat" ]; then
    echo "Using $CONTRACTS_DIR as HH proj"
    cd $CONTRACTS_DIR && npx hardhat node --no-deploy --hostname $L1_HOST --port $L1_PORT &
    L1_PID=$!
    PIDS+=$L1_PID
    echo "L1 PID: $L1_PID"
    sleep 3
else
    echo "invalid value for L1_STACK: $L1_STACK"
    exit 1
fi

# Optionally deploy the contracts
if [ "$L1_DEPLOY" = "true" ]; then
    echo "Deploying contracts..."
    bash $SBIN/deploy_l1_contracts.sh
fi

# Follow output
# tail -f $L1_PID
echo "L1 started... (Use ctrl-c to stop)"
wait $L1_PID

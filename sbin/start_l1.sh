#!/bin/sh
SBIN=`dirname $0`
SBIN="`cd "$SBIN"; pwd`"
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
ENV=".genesis.env"
if ! test -f $ENV; then
    echo "Expected dotenv at $ENV (does not exist)."
    exit
fi
echo "Using dotenv: $ENV"
. $ENV

L1_HOST=`echo $L1_ENDPOINT | awk -F':' '{print substr($2, 3)}'`
L1_WS_PORT=`echo $L1_ENDPOINT | awk -F':' '{print $3}'`
echo "Parsed endpoint ($L1_HOST) and port: $L1_WS_PORT from $L1_ENDPOINT"

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
      --ws.port $L1_WS_PORT &

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
    cd $CONTRACTS_DIR && npx hardhat node --no-deploy --hostname $L1_HOST --port $L1_WS_PORT &
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
if [ "$L1_STACK" = "geth" ]; then
    tail -f $L1_PID
elif [ "$L1_STACK" = "hardhat" ]; then
    tail -f $L1_PID
fi

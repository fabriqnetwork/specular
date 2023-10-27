#!/bin/bash
if [ ! -d "$CONFIG_DIR" ]; then
    SBIN=`dirname $0`
    SBIN="`cd "$SBIN"; pwd`"
    . $SBIN/configure.sh
    CONFIG_DIR="`cd "$CONFIG_DIR"; pwd`"
fi
echo "Using $CONFIG_DIR as HH proj"

# Check that the dotenv exists.
ENV=".genesis.env"
if ! test -f $ENV; then
    echo "Expected dotenv at $ENV (does not exist)."
    exit
fi
echo "Using dotenv: $ENV"
. $ENV

# Define a function to convert a path to be relative to another directory.
relpath () {
    echo `python3 -c "import os.path; print(os.path.relpath('$1', '$2'))"`
}

# Get relative paths, since we have to run `create_genesis.ts` from the HH proj.
BASE_GENESIS_PATH=`relpath $BASE_GENESIS_PATH $CONFIG_DIR`
GENESIS_PATH=`relpath $GENESIS_PATH $CONFIG_DIR`

L1_WS_PORT=`echo $L1_ENDPOINT | awk -F':' '{print $3}'`
echo "Parsed port: $L1_WS_PORT from $L1_ENDPOINT"

# Start an L1 if one isn't already running
if ! ss -tuln | grep -q ":$L1_WS_PORT "; then
    echo "Starting L1..."
    if [ $L1_STACK = "geth" ]; then
        $GETH_BIN \
        --dev \
        --dev.period 1 \
        --verbosity 3 \
        --http \
        --http.api eth,web3,net \
        --http.addr 0.0.0.0 \
        --ws \
        --ws.api eth,net,web3 \
        --ws.addr 0.0.0.0 \
        --ws.port $L1_WS_PORT > /dev/null 2>&1 &
        sleep 3
        L1_PID=$!

    elif [ $L1_STACK = "hardhat" ]; then
        cd $CONTRACTS_DIR && npx hardhat node --no-deploy --port $L1_WS_PORT &
        L1_PID=$!
        echo "L1 PID: $L1_PID"
        sleep 3
    else
        echo "invalid value for L1_STACK: $L1_STACK"
        exit 1
    fi
else
    L1_PID=0
    echo "L1 already started..."
fi


# Create genesis.json file.
cd $CONFIG_DIR && npx ts-node src/create_genesis.ts --in $BASE_GENESIS_PATH --out $GENESIS_PATH --l1network "$L1_ENDPOINT"

if [ $L1_PID -ne 0 ]; then
    echo "Stopping L1"
    kill $L1_PID
fi

# If the contracts directory exists, initialize a reference to the genesis file at
# "contracts/.genesis" (using relative paths as appropriate).
if [ -d "$CONTRACTS_DIR" ]; then
    CONTRACTS_DIR=`cd $CONTRACTS_DIR; pwd`
    CONTRACTS_ENV=$CONTRACTS_DIR/$ENV
    # If it already exists, check if we should overwrite the file.
    if test -f $CONTRACTS_ENV; then
        read -r -p "Overwrite $CONTRACTS_ENV with a new file? [y/N] " response
        if [[ "$response" =~ ^([yY][eE][sS]|[yY])$ ]]; then
            rm $CONTRACTS_ENV
        else
            exit
        fi
    fi
    # Write file, using relative paths.
    echo "Initializing $CONTRACTS_ENV"
    GENESIS_PATH=`relpath $GENESIS_PATH $CONTRACTS_DIR`
    echo GENESIS_PATH=$GENESIS_PATH >> $CONTRACTS_ENV
fi

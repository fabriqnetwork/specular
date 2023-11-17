#!/bin/bash
# Check that the dotenv exists.
ENV=".genesis.env"
if ! test -f $ENV; then
    echo "Expected dotenv at $ENV (does not exist)."
    exit
fi
echo "Using dotenv: $ENV"
. $ENV

if [ ! -d "$OPS_DIR" ]; then
    SBIN=`dirname $0`
    SBIN="`cd "$SBIN"; pwd`"
    . $SBIN/configure.sh
    OPS_DIR="`cd "$OPS_DIR"; pwd`"
fi
echo "Using $OPS_DIR as ops directory."

# Define a function to convert a path to be relative to another directory.
relpath () {
    echo `python3 -c "import os.path; print(os.path.relpath('$1', '$2'))"`
}

# Define a function that requests a user to confirm
# that overwriting file ($1) is okay, if it exists.
guard_overwrite () {
    if test -f $1; then
	read -r -p "Overwrite $1 with a new file? [y/N] " response
	if [[ "$response" =~ ^([yY][eE][sS]|[yY])$ ]]; then
	    rm $1
	else
	    exit
	fi
    fi
}

# Get relative paths
GENESIS_PATH=`relpath $GENESIS_PATH $OPS_DIR`

# Create genesis.json file.
echo "Generating new genesis file at $GENESIS_PATH"
cd $OPS_DIR
guard_overwrite $GENESIS_PATH\
go run ./cmd/genesis/main.go \
    --genesis-config $GENESIS_CFG_PATH \
    --out $GENESIS_PATH \
    --l1-rpc-url $L1_ENDPOINT \
    --export-hash $GENESIS_EXPORTED_HASH_PATH

# Initialize a reference to the genesis file at
# "contracts/.genesis" (using relative paths as appropriate).
CONTRACTS_ENV=$CONTRACTS_DIR/$ENV
guard_overwrite $CONTRACTS_ENV
# Write file, using relative paths.
echo "Initializing $CONTRACTS_ENV"
GENESIS_PATH=`relpath $GENESIS_PATH $CONTRACTS_DIR`
echo GENESIS_PATH=$GENESIS_PATH >> $CONTRACTS_ENV

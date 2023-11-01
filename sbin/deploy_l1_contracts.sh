#!/bin/bash
if [ ! -d "$CONTRACTS_DIR" ]; then
    SBIN=`dirname $0`
    SBIN="`cd "$SBIN"; pwd`"
    . $SBIN/configure.sh
    CONTRACTS_DIR="`cd "$CONTRACTS_DIR"; pwd`"
fi
echo "Using $CONTRACTS_DIR as HH proj"

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

# Define a function that requests a user to confirm
# that overwriting file ($1) is okay, if it exists.
guard_overwrite () {
    if test -f $1; then
	read -r -p "Overwrite $1 with a new file? [y/N] " response
	if [[ "$response" =~ ^([yY][eE][sS]|[yY])$ ]]; then
	    rm $1
	else
	    exit 1
	fi
    fi
}

# Get relative paths, since we have to run `create_genesis.ts` from the HH proj.
BASE_ROLLUP_CFG_PATH=`relpath $BASE_ROLLUP_CFG_PATH $CONTRACTS_DIR`
ROLLUP_CFG_PATH=`relpath $ROLLUP_CFG_PATH $CONTRACTS_DIR`
GENESIS_PATH=`relpath $GENESIS_PATH $CONTRACTS_DIR`

# Generate genesis file
$SBIN/create_genesis.sh

cd $CONTRACTS_DIR
echo "Deploying l1 contracts..."
npx hardhat deploy --network localhost

echo "Generating rollup config..."
guard_overwrite $ROLLUP_CFG_PATH
npx ts-node scripts/config/create_config.ts \
  --in $BASE_ROLLUP_CFG_PATH \
  --out $ROLLUP_CFG_PATH \
  --genesis $GENESIS_PATH \
  --l1-network $L1_ENDPOINT

echo "Done."

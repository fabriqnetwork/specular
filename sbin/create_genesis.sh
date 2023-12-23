#!/bin/bash
# Currently the local sbin paths are relative to the project root.
SBIN=$(dirname "$(readlink -f "$0")")
SBIN="$(
  cd "$SBIN"
  pwd
)"
ROOT_DIR=$SBIN/..

# Check that the all required dotenv files exists.
PATHS_ENV=".paths.env"
if ! test -f "$PATHS_ENV"; then
  echo "Expected dotenv at $PATHS_ENV (does not exist)."
  exit
fi
echo "Using paths dotenv: $PATHS_ENV"
. $PATHS_ENV

GENESIS_ENV=".genesis.env"
if ! test -f "$GENESIS_ENV"; then
  echo "Expected dotenv at $GENESIS_ENV (does not exist)."
  exit
fi
echo "Using genesis dotenv: $GENESIS_ENV"
. $GENESIS_ENV

echo "Using $OPS_DIR as ops directory."

# Define a function to convert a path to be relative to another directory.
relpath() {
  echo $(python3 -c "import os.path; print(os.path.relpath('$1', '$2'))")
}

# Define a function that requests a user to confirm
# that overwriting file ($1) is okay, if it exists.
guard_overwrite() {
  if test -f $1; then
    read -r -p "Overwrite $1 with a new file? [y/N] " response
    if [[ "$response" =~ ^([yY][eE][sS]|[yY])$ ]]; then
      rm $1
    else
      exit 1
    fi
  fi
}

# Get relative paths for $OPS_DIR
GENESIS_CFG_PATH=$(relpath $GENESIS_CFG_PATH $OPS_DIR)
GENESIS_PATH=$(relpath $GENESIS_PATH $OPS_DIR)
GENESIS_EXPORTED_HASH_PATH=$(relpath $GENESIS_EXPORTED_HASH_PATH $OPS_DIR)
echo "Generating new genesis file at $GENESIS_PATH and exporting hash to $GENESIS_EXPORTED_HASH_PATH"
cd $OPS_DIR
guard_overwrite $GENESIS_PATH
# Create genesis.json file.
CMD="""
$OPS_GENESIS_BIN \
    --genesis-config $GENESIS_CFG_PATH \
    --out $GENESIS_PATH \
    --l1-rpc-url $L1_ENDPOINT \
    --export-hash $GENESIS_EXPORTED_HASH_PATH
"""
echo "Running $CMD"
eval $CMD

# Initialize a reference to the genesis file at
# "contracts/.genesis" (using relative paths as appropriate).
CONTRACTS_ENV=$CONTRACTS_DIR/$GENESIS_ENV
guard_overwrite $CONTRACTS_ENV
# Write file, using relative paths.
echo "Initializing contracts dotenv $CONTRACTS_ENV"
GENESIS_PATH=$(relpath $GENESIS_PATH $CONTRACTS_DIR)
GENESIS_EXPORTED_HASH_PATH=$(relpath $GENESIS_EXPORTED_HASH_PATH $CONTRACTS_DIR)
echo GENESIS_PATH=$GENESIS_PATH >>$CONTRACTS_ENV
echo GENESIS_EXPORTED_HASH_PATH=$GENESIS_EXPORTED_HASH_PATH >>$CONTRACTS_ENV

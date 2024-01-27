#!/bin/bash

set -e

# Check that a dotenv exists.
geth_env=".sp_geth.env"
if ! test -f "$geth_env"; then
  echo "Error: expected dotenv at ./$geth_env (does not exist); could not clean current working directory."
  exit
fi
source $geth_env
echo "Removing sp-geth data directory: $DATA_DIR"
rm -rf $DATA_DIR

if test -f .start_sp_geth.sh.lock; then
  echo "Removing docker lock file"
  lock_file=$WAIT_DIR/.start_sp_geth.sh.lock
fi

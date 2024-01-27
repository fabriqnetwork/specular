#!/bin/bash

set -e

# Check that a dotenv exists.
dotenv_file=".sp_geth.env"
if ! test -f "$dotenv_file"; then
  echo "Error: expected dotenv at ./$dotenv_file (does not exist); could not clean current working directory."
  exit
fi
source $dotenv_file
echo "Removing sp-geth data directory: $DATA_DIR"
rm -rf $DATA_DIR

if test -f .start_sp_geth.sh.lock; then
  lock_file=".start_sp_geth.sh.lock"
  echo "Removing docker lock file"
fi

#!/bin/bash
# Check that a dotenv exists.
GETH_ENV=".sp_geth.env"
if ! test -f "$GETH_ENV"; then
  echo "expected dotenv at ./$GETH_ENV (does not exist); could not clean cwd."
  exit
fi
. $GETH_ENV
echo "Removing sp-geth data dir $DATA_DIR"
rm -rf $DATA_DIR

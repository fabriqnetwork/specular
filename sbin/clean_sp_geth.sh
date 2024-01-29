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

if test -f .sp_geth_started.lock; then

  echo "Removing docker lock file"
  L1_WAIT=$WAIT_DIR/.sp_geth_started.lock
fi

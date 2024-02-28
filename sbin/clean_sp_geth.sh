#!/bin/bash
# Check that a dotenv exists.
WORKSPACE_DIR=$HOME/.spc/workspaces/active_workspace
GETH_ENV="$WORKSPACE_DIR/.sp_geth.env"

if ! test -f "$GETH_ENV"; then
  echo "expected dotenv at $GETH_ENV (does not exist); could not clean cwd."
  exit
fi
. $GETH_ENV
echo "Removing sp-geth data dir $DATA_DIR"
rm -rf $DATA_DIR

if test -f .start_sp_geth.sh.lock; then

  echo "Removing docker lock file"
  L1_WAIT=$WAIT_DIR/.start_sp_geth.sh.lock
fi

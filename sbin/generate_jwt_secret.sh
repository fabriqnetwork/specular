#!/bin/bash
SBIN=$(dirname "$(readlink -f "$0")")
SBIN="$(
  cd "$SBIN"
  pwd
)"
. $SBIN/utils/utils.sh
. $SBIN/utils/crypto.sh

WORKSPACE_DIR=$HOME/.spc/workspaces/active_workspace

ROOT_DIR=$SBIN/..
JWT=$(generate_jwt_secret)

# Write to sp-magi's expected JWT secret path.
SP_MAGI_ENV=$WORKSPACE_DIR/.sp_magi.env
reqdotenv "sp_magi" $SP_MAGI_ENV
echo $JWT >$JWT_SECRET_PATH

# Write to sp-geth's expected JWT secret path.
SP_GETH_ENV=$WORKSPACE_DIR/.sp_geth.env
reqdotenv "sp_geth" $SP_GETH_ENV
echo $JWT >$JWT_SECRET_PATH

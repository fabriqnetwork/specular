#!/bin/bash
SBIN=$(dirname "$(readlink -f "$0")")
SBIN="$(
  cd "$SBIN"
  pwd
)"
. $SBIN/utils/utils.sh
. $SBIN/utils/crypto.sh

ROOT_DIR=$SBIN/..
JWT=$(generate_jwt_secret)
# Write to sp-magi's expected JWT secret path.
reqdotenv "sp_magi" ".sp_magi.env"
echo $JWT >$JWT_SECRET_PATH
# Write to sp-geth's expected JWT secret path.
reqdotenv "sp_magi" ".sp_geth.env"
echo $JWT >$JWT_SECRET_PATH

#!/bin/bash

set -e

# Get the directory of the script
SCRIPT_DIR=$(dirname "$(readlink -f "$0")")
SCRIPT_DIR="$(
  cd "$SCRIPT_DIR"
  pwd
)"
# Sourcing utility scripts
. $SCRIPT_DIR/utils/utils.sh
. $SCRIPT_DIR/utils/crypto.sh

# Setting the root directory
ROOT_DIR=$SCRIPT_DIR/..

# Generating JWT secret
JWT_SECRET=$(generate_jwt_secret)

# Write JWT secret to sp-magi's expected path
require_dotenv "sp_magi" ".sp_magi.env"
echo $JWT_SECRET >$JWT_SECRET_PATH

# Write JWT secret to sp-geth's expected path
require_dotenv "sp_magi" ".sp_geth.env"
echo $JWT_SECRET >$JWT_SECRET_PATH

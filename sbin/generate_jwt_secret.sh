#!/bin/bash

set -e

# Get the directory of the script
script_dir=$(dirname "$(readlink -f "$0")")
script_dir="$(
  cd "$script_dir"
  pwd
)"
# Source utility scripts
. $script_dir/utils/utils.sh
. $script_dir/utils/crypto.sh

# Set the root directory
root_dir=$script_dir/..

# Generate JWT secret
jwt_secret=$(generate_jwt_secret)

# Write JWT secret to sp-magi's expected path
require_dotenv "sp_magi" ".sp_magi.env"
echo $jwt_secret >$JWT_SECRET_PATH

# Write JWT secret to sp-geth's expected path
require_dotenv "sp_geth" ".sp_geth.env"
echo $jwt_secret >$JWT_SECRET_PATH

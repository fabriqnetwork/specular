#!/bin/bash

set -e

# Get the directory of the script
script_dir=$(dirname $0)

# Clean up Ethereum related scripts
$script_dir/clean_sp_geth.sh
$script_dir/clean_deployment.sh

# Remove dotenv files
rm -f .contracts.env
rm -f .genesis.env
rm -f .sp_geth.env
rm -f .sp_magi.env
rm -f .sidecar.env
rm -f .paths.env
echo "Done removing dotenv files."

# Remove JWT secret file
echo "Removing $JWT_SECRET_PATH"
rm -f $JWT_SECRET_PATH

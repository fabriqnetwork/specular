#!/bin/bash

set -e

# Change directory to the workspace
cd /specular/workspace

# Set environment variables
export INFURA_KEY=$(cat infura_private_key.txt)
export DEPLOYER_PRIVATE_KEY=$(cat deployer_private_key.txt)
export SEQUENCER_PRIVATE_KEY=$(cat sequencer_private_key.txt)
export VALIDATOR_PRIVATE_KEY=$(cat validator_private_key.txt)

# Source environment files
set -o allexport
. .sp_geth.env
. .sp_magi.env
. .contracts.env
. .genesis.env
. .paths.env
. .sidecar.env
set +o allexport

# Handle command line arguments
case "$1" in
deploy)
    # Run the main container command for deployment
    /specular/sbin/generate_jwt_secret.sh
    /specular/sbin/deploy_l1_contracts.sh -y
    ;;
start)
    shift
    /specular/sbin/$@
    ;;
*)
    echo "Unknown Command"
    exit 1
    ;;
esac

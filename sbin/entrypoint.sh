#!/bin/bash
set -ex
cd /specular/workspace

# remove all locks since this need to run first
# LOCKFILES=`ls .*.lock 2> /dev/null | wc -l`
# if [[ $LOCKFILES -gt 0 ]]; then
#     echo "Removing lockfiles"
#     rm .*.lock
# fi

echo "Setting environment variables"
export INFURA_KEY=$(cat infura_pk.txt)
export DEPLOYER_PRIVATE_KEY=$(cat deployer_pk.txt)
export SEQUENCER_PRIVATE_KEY=$(cat sequencer_pk.txt)
export VALIDATOR_PRIV_KEY=$(cat validator_pk.txt)
set -o allexport
. .sp_geth.env
. .sp_magi.env
. .contracts.env
. .genesis.env
. .paths.env
. .sidecar.env
set +o allexport

case "$1" in
deploy)
  # Run the main container command.
  echo "Running deploy for genesis and JWT"
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

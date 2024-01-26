#!/bin/bash
set -ex
cd /specular/workspace

# remove all locks since this need to run first
LOCKFILES=`ls .*.lock 2> /dev/null | wc -l`
if [[ $LOCKFILES -gt 0 ]]; then
    echo "Removing lockfiles"
    rm .*.lock
fi


echo "Setting environment variables"
export INFURA_KEY=`cat infura_pk.txt`
export DEPLOYER_PRIVATE_KEY=`cat deployer_pk.txt`
export SEQUENCER_PRIVATE_KEY=`cat sequencer_pk.txt`
export VALIDATOR_PRIV_KEY=`cat validator_pk.txt`

printenv

case "$1" in
deploy)
    # Run the main container command.
    echo "Running deploy for genesis and JWT"
    /bin/bash ../sbin/generate_jwt_secret.sh
    /bin/bash ../sbin/deploy_l1_contracts.sh -y
    ;;
start)
    shift
    /bin/bash ../sbin/$@
    ;;
*)
  echo "Unknown Command"
  exit 1
  ;;
esac


#!/bin/bash

# the local sbin paths are relative to the project root
SBIN=$(dirname "$(readlink -f "$0")")
SBIN="$(
  cd "$SBIN"
  pwd
)"
. $SBIN/utils/utils.sh
ROOT_DIR=$SBIN/..

reqdotenv "paths" ".paths.env"
reqdotenv "deployments" ".deployments.env"

cd $CONTRACTS_DIR
npx hardhat run scripts/pause.ts

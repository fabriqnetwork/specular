#!/bin/bash
set -e

# currently the local sbin paths are relative to the project root
SBIN=$(dirname "$(readlink -f "$0")")
SBIN="$(
  cd "$SBIN"
  pwd
)"
. $SBIN/utils/utils.sh
ROOT_DIR=$SBIN/..

WAITFILE="/tmp/.${0##*/}.lock"

if [[ ! -z ${WAIT_DIR+x} ]]; then
  WAITFILE=$WAIT_DIR/.${0##*/}.lock
fi

WORKSPACE_DIR=$HOME/.spc/workspaces/active_workspace

PATHS_ENV=$WORKSPACE_DIR/.paths.env
SP_GETH_ENV=$WORKSPACE_DIR/.sp_geth.env

reqdotenv "paths" $PATHS_ENV
reqdotenv "sp_geth" $SP_GETH_ENV

# Parse args.
optspec="chw"
while getopts "$optspec" optchar; do
  case "${optchar}" in
  w)
    WAIT=true
    ;;
  c)
    echo "Cleaning..."
    $SBIN/clean_sp_geth.sh
    ;;
  h)
    echo "usage: $0 [-c][-h][-w]"
    echo "-c : clean before running"
    echo "-w : generate docker-compose wait for file"
    exit
    ;;
  *)
    if [ "$OPTERR" != 1 ] || [ "${optspec:0:1}" = ":" ]; then
      echo "Unknown option: '-${OPTARG}'"
      exit 1
    fi
    ;;
  esac
done

if [ "$WAIT" = "true" ]; then
  if test -f $WAITFILE; then
    echo "Removing wait file for docker..."
    rm $WAITFILE
  fi
fi

spc up spgeth

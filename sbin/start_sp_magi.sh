#!/bin/bash
set -e

# the local sbin paths are relative to the project root
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
SP_MAGI_ENV=$WORKSPACE_DIR/.sp_magi.env

reqdotenv "paths" $PATHS_ENV
reqdotenv "sp_magi" $SP_MAGI_ENV

# Generate waitfile for service init (docker/k8)
WAITFILE="/tmp/.${0##*/}.lock"

if [[ ! -z ${WAIT_DIR+x} ]]; then
  WAITFILE=$WAIT_DIR/.${0##*/}.lock
fi

# Set sync flags.
SYNC_FLAGS=""
if [ $SYNC_MODE = "checkpoint" ]; then
  echo "Enabling checkpoint."
  SYNC_FLAGS="--checkpoint"
fi

# Set devnet flags.
DEVNET_FLAGS=""
if [ "$DEVNET" = true ]; then
  echo "Enabling devnet mode."
  DEVNET_FLAGS="--devnet"
fi

# Set local sequencer flags.
SEQUENCER_FLAGS=""
if [ "$SEQUENCER" = true ]; then
  echo "Enabling local sequencer."
  SEQUENCER_FLAGS="--sequencer"
fi

spc up spmagi $SYNC_FLAGS $DEVNET_FLAGS $SEQUENCER_FLAGS

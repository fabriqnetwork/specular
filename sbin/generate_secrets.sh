#!/bin/bash
optspec="djyw"
NUM_ACCOUNTS=0
AUTO_ACCEPT=false
while getopts "$optspec" optchar; do
  case "${optchar}" in
  y)
    AUTO_ACCEPT=true
    ;;
  d)
    GEN_DEPLOYER=true
    ;;
  j)
    GEN_JWT=true
    ;;
  w)
    WAIT=true
    ;;
  *)
    echo "usage: $0 [-d][-j][-y][-h][-w]"
    echo "-d : generate deployer"
    echo "-j : generate jwt secret"
    echo "-y : auto accept prompts"
    echo "-w : generate docker-compose wait for file"
    exit
    ;;
  esac
done

# the local sbin paths are relative to the project root
SBIN=$(dirname "$(readlink -f "$0")")
SBIN="$(
  cd "$SBIN"
  pwd
)"
. $SBIN/utils/utils.sh
. $SBIN/utils/crypto.sh
ROOT_DIR=$SBIN/..

reqdotenv "sp_magi" ".sp_magi.env"
reqdotenv "sidecar" ".sidecar.env"
reqdotenv "paths" ".paths.env"

# Generate waitfile for service init (docker/k8)
WAITFILE="/tmp/.${0##*/}.lock"

if [[ ! -z ${WAIT_DIR+x} ]]; then
  WAITFILE=$WAIT_DIR/.${0##*/}.lock
fi

if [ "$WAIT" = "true" ]; then
  if test -f $WAITFILE; then
    echo "Removing wait file for docker..."
    rm $WAITFILE
  fi
fi

CONTRACTS_ENV=".contracts.env"
guard_overwrite $CONTRACTS_ENV $AUTO_ACCEPT

# Generate accounts
VALIDATOR_ADDRESS=$(generate_wallet $VALIDATOR_PK_PATH)
echo "Generated account (address=$VALIDATOR_ADDRESS, priv_key_path=$VALIDATOR_PK_PATH)"
SEQUENCER_ADDRESS=$(generate_wallet $SEQUENCER_PK_FILE)
echo "Generated account (address=$SEQUENCER_ADDRESS, priv_key_path=$SEQUENCER_PK_FILE)"
if [ "$DISSEMINATOR_PK_PATH" != "$SEQUENCER_PK_FILE" ]; then
  echo "$DISSEMINATOR_PK_PATH" "$SEQUENCER_PK_FILE"
  guard_overwrite $DISSEMINATOR_PK_PATH $AUTO_ACCEPT
  cat $SEQUENCER_PK_FILE >$DISSEMINATOR_PK_PATH
fi

# Write dotenv
echo "VALIDATOR_ADDRESS=$VALIDATOR_ADDRESS" >$CONTRACTS_ENV
echo "SEQUENCER_ADDRESS=$SEQUENCER_ADDRESS" >>$CONTRACTS_ENV
echo "Wrote addresses to $CONTRACTS_ENV"

# Generate deployer account
if [ "$GEN_DEPLOYER" = "true" ]; then
  deployer_pk_path=deployer_pk.txt
  DEPLOYER_ADDRESS=$(generate_wallet $deployer_pk_path)
  echo "Generated account (address=$DEPLOYER_ADDRESS, priv_key_path=$deployer_pk_path)"
  echo "DEPLOYER_ADDRESS=$DEPLOYER_ADDRESS" >>$CONTRACTS_ENV
  echo "DEPLOYER_PRIVATE_KEY=$(cat $deployer_pk_path)" >>$CONTRACTS_ENV
  echo "Wrote address to $CONTRACTS_ENV"
fi

if [ "$GEN_JWT" = "true" ]; then
  JWT=$(generate_jwt_secret)
  echo $JWT >"$JWT_SECRET_PATH"
fi

if [ "$WAIT" = "true" ]; then
  echo "Creating wait file for docker at $WAITFILE..."
  touch $WAITFILE
fi

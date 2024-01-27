#!/bin/bash

set -e

optspec="djyw"
num_accounts=0
auto_accept=false
while getopts "$optspec" optchar; do
  case "${optchar}" in
  y)
    auto_accept=true
    ;;
  d)
    gen_deployer=true
    ;;
  j)
    gen_jwt=true
    ;;
  w)
    wait_flag=true
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
sbin=$(dirname "$(readlink -f "$0")")
sbin="$(
  cd "$sbin"
  pwd
)"
. $sbin/utils/utils.sh
. $sbin/utils/crypto.sh
root_dir=$sbin/..

require_dotenv "sp_magi" ".sp_magi.env"
require_dotenv "sidecar" ".sidecar.env"
require_dotenv "paths" ".paths.env"

# Generate waitfile for service init (docker/k8)
waitfile="/tmp/.${0##*/}.lock"

if [[ ! -z ${WAIT_DIR+x} ]]; then
  waitfile=$WAIT_DIR/.${0##*/}.lock
fi

if [ "$wait_flag" = "true" ]; then
  if test -f $waitfile; then
    echo "Removing wait file for docker..."
    rm $waitfile
  fi
fi

contracts_env=".contracts.env"
confirm_overwrite $contracts_env $auto_accept

# Generate accounts
validator_address=$(generate_wallet $validator_pk_path)
echo "Generated account (address=$validator_address, priv_key_path=$validator_pk_path)"
sequencer_address=$(generate_wallet $sequencer_pk_file)
echo "Generated account (address=$sequencer_address, priv_key_path=$sequencer_pk_file)"
if [ "$disseminator_pk_path" != "$sequencer_pk_file" ]; then
  confirm_overwrite $disseminator_pk_path $auto_accept
  cat $sequencer_pk_file >$disseminator_pk_path
fi

# Write dotenv
echo "VALIDATOR_ADDRESS=$validator_address" >$contracts_env
echo "SEQUENCER_ADDRESS=$sequencer_address" >>$contracts_env
echo "Wrote addresses to $contracts_env"

# Generate deployer account
if [ "$gen_deployer" = "true" ]; then
  deployer_pk_path=deployer_pk.txt
  deployer_address=$(generate_wallet $deployer_pk_path)
  echo "Generated account (address=$deployer_address, priv_key_path=$deployer_pk_path)"
  echo "DEPLOYER_ADDRESS=$deployer_address" >>$contracts_env
  echo "DEPLOYER_PRIVATE_KEY=$(cat $deployer_pk_path)" >>$contracts_env
  echo "Wrote address to $contracts_env"
fi

if [ "$gen_jwt" = "true" ]; then
  jwt=$(generate_jwt_secret)
  echo $jwt >"$jwt_secret_path"
fi

if [ "$wait_flag" = "true" ]; then
  echo "Creating wait file for docker at $waitfile..."
  touch $waitfile
fi

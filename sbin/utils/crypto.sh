#!/bin/bash

set -e

# Generates a wallet.
# Writes the address to stdout and private key to $1
generate_wallet() {
  local wallet_info
  wallet_info=$(cast wallet new)
  local address
  address=$(echo "$wallet_info" | awk '/Address/ { print $2 }')
  local private_key
  private_key=$(echo "$wallet_info" | awk '/Private key/ { print $3 }')
  guard_overwrite $1
  echo -n "$private_key" >"$1"
  echo "$address"
}

generate_jwt_secret() {
  openssl rand -hex 32
}

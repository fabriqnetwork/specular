#!/bin/bash

# Generates a wallet.
# Writes the address to stdout and private key to $1
generate_wallet() {
  wallet=$(cast wallet new)
  address=$(echo "$wallet" | awk '/Address/ { print $2 }')
  priv_key=$(echo "$wallet" | awk '/Private key/ { print $3 }')
  echo "$address"
  guard_overwrite $1
  echo $priv_key | tr -d '\n' >$1
}

generate_jwt_secret() {
  openssl rand -hex 32
}

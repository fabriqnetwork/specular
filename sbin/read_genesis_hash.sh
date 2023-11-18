#!/bin/bash

# Check that the dotenv exists
ENV=".genesis.env"
if ! test -f $ENV ; then
    echo "Expected dotenv at $ENV (does not exist)."
    exit 1
fi
. $ENV

if ! test -f "$GENESIS_EXPORTED_HASH_PATH" ; then
    echo "Expected GENESIS_EXPORTED_HASH_PATH to be set in $ENV."
    exit 1
fi

echo "$(cat $GENESIS_EXPORTED_HASH_PATH)"

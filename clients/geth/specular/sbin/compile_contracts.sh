#!/bin/bash
SBIN=`dirname $0`
SBIN="`cd "$SBIN"; pwd`"
. $SBIN/configure.sh
cd $CONTRACTS_DIR && npm install --force && npx hardhat compile

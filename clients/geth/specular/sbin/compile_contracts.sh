#!/bin/bash
SBIN=`dirname $0`
SBIN="`cd "$SBIN"; pwd`"
. $SBIN/configure.sh
cd $CONTRACTS_DIR && yarn install && npx hardhat compile

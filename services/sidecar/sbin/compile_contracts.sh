#!/bin/sh
SBIN=`dirname $0`
SBIN="`cd "$SBIN"; pwd`"
. $SBIN/configure.sh
cd $CONTRACTS_DIR && npx hardhat compile

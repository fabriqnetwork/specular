#!/bin/sh
SBIN=`dirname $0`
SBIN="`cd "$SBIN"; pwd`"
. $SBIN/configure.sh
cd $CONTRACTS_DIR && pnpm install && npx hardhat compile

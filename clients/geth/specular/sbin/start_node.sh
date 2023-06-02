#!/bin/bash
SBIN=`dirname $0`
SBIN="`cd "$SBIN"; pwd`"
. $SBIN/configure.sh

cd $DATA_DIR && $GETH_SPECULAR_DIR/build/bin/geth "$@"

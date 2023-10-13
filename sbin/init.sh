#!/bin/bash
SBIN=`dirname $0`
SBIN="`cd "$SBIN"; pwd`"

. $SBIN/configure.sh
. $SBIN/configure_system.sh

cd $DATA_DIR
$GETH_DIR/build/bin/geth \
  --datadir ./ \
  --networkid $NETWORK_ID \
  init ./genesis.json

if [[ $USE_CLEF == 'true' ]]; then
  cd $SBIN
  ./init_clef.exp
fi



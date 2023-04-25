#!/bin/bash
SBIN=`dirname $0`
SBIN="`cd "$SBIN"; pwd`"

. $SBIN/configure.sh
. $SBIN/configure_system.sh

cd $DATA_DIR
$GETH_SPECULAR_DIR/build/bin/geth \
  --datadir ./data_sequencer \
  --networkid $NETWORK_ID \
  init ./genesis.json

$GETH_SPECULAR_DIR/build/bin/geth \
  --datadir ./data_validator \
  --networkid $NETWORK_ID \
  init ./genesis.json

$GETH_SPECULAR_DIR/build/bin/geth \
  --datadir ./data_indexer \
  --networkid $NETWORK_ID \
  init ./genesis.json

if [[ $USE_CLEF == 'true' ]]; then
  cd $SBIN
  ./init_clef.exp
fi



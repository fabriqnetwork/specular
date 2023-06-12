#!/bin/bash
SBIN=`dirname $0`
SBIN="`cd "$SBIN"; pwd`"

. $SBIN/configure.sh
. $SBIN/configure_geth_args.sh # sets NETWORK_ID

cd $DATA_DIR
$SPECULAR_CLIENT_DIR/build/bin/geth \
  --datadir ./data_sequencer \
  --networkid $NETWORK_ID \
  init ./genesis.json

$SPECULAR_CLIENT_DIR/build/bin/geth \
  --datadir ./data_validator \
  --networkid $NETWORK_ID \
  init ./genesis.json

$SPECULAR_CLIENT_DIR/build/bin/geth \
  --datadir ./data_indexer \
  --networkid $NETWORK_ID \
  init ./genesis.json

if [ $USE_CLEF -eq "1" ]; then
  cd $SBIN
  ./init_clef.exp
fi

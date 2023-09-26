#!/bin/bash
SBIN=`dirname $0`
SBIN="`cd "$SBIN"; pwd`"
. $SBIN/configure.sh
. $SBIN/configure_system.sh
cd $DATA_DIR

CLEF_PW="unsafe-password"
echo ${CLEF_PW} | \
$SIDECAR_DIR/build/bin/clef \
    --suppress-bootwarn \
    --configdir ./data_clef/ \
    --auditlog ./data_clef/audit.log \
    --signersecret ./data_clef/masterseed.json \
    --keystore ./data_sequencer/keystore/ \
    --rules ./ruleset.js \
    --http \
    --advanced \
    --chainid $L1_CHAIN_ID

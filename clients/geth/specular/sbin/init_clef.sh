#!/bin/bash
SBIN=`dirname $0`
SBIN="`cd "$SBIN"; pwd`"
. $SBIN/configure.sh
. $SBIN/configure_system.sh
cd $DATA_DIR

rm -rf ./data_clef/

$GETH_SPECULAR_DIR/build/bin/clef init \
    --configdir ./data_clef

$GETH_SPECULAR_DIR/build/bin/clef setpw \
    --configdir ./data_clef/ \
    $SEQUENCER_ADDR

$GETH_SPECULAR_DIR/build/bin/clef setpw \
    --configdir ./data_clef/ \
    $VALIDATOR_ADDR

$GETH_SPECULAR_DIR/build/bin/clef setpw \
    --configdir ./data_clef/ \
    $INDEXER_ADDR

CHECKSUM=$(shasum -a 256 ./ruleset.js)
ARR=($CHECKSUM)
echo "rule checksum is: ${ARR[0]}"
$GETH_SPECULAR_DIR/build/bin/clef attest \
    --configdir ./data_clef \
    --signersecret ./data_clef/masterseed.json \
    ${ARR[0]}

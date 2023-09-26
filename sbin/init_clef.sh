#!/bin/bash
SBIN=`dirname $0`
SBIN="`cd "$SBIN"; pwd`"
. $SBIN/configure.sh
. $SBIN/configure_system.sh

cd $DATA_DIR

rm -rf ./data_clef/

$SIDECAR_DIR/build/bin/clef init \
    --configdir ./data_clef

$SIDECAR_DIR/build/bin/clef setpw \
    --configdir ./data_clef/ \
    $SEQUENCER_ADDR

$SIDECAR_DIR/build/bin/clef setpw \
    --configdir ./data_clef/ \
    $VALIDATOR_ADDR

$SIDECAR_DIR/build/bin/clef setpw \
    --configdir ./data_clef/ \
    $INDEXER_ADDR

CHECKSUM=$(shasum -a 256 ./ruleset.js)
ARR=($CHECKSUM)
echo "rule checksum is: ${ARR[0]}"
$SIDECAR_DIR/build/bin/clef attest \
    --configdir ./data_clef \
    --signersecret ./data_clef/masterseed.json \
    ${ARR[0]}

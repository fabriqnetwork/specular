#!/bin/bash
SBIN=`dirname $0`
SBIN="`cd "$SBIN"; pwd`"
. $SBIN/configure.sh
cd $DATA_DIR

rm -rf ./data_clef/

$GETH_SPECULAR_DIR/build/bin/clef init \
    --configdir ./data_clef

$GETH_SPECULAR_DIR/build/bin/clef setpw \
    --configdir ./data_clef/ \
    0xf39fd6e51aad88f6f4ce6ab8827279cfffb92266

$GETH_SPECULAR_DIR/build/bin/clef setpw \
    --configdir ./data_clef/ \
    0x70997970c51812dc3a010c7d01b50e0d17dc79c8

$GETH_SPECULAR_DIR/build/bin/clef setpw \
    --configdir ./data_clef/ \
    0x3c44cdddb6a900fa2b585dd299e03d12fa4293bc

CHECKSUM=$(shasum -a 256 ./ruleset.js)
ARR=($CHECKSUM)
echo "rule checksum is: ${ARR[0]}"
$GETH_SPECULAR_DIR/build/bin/clef attest \
    --configdir ./data_clef \
    --signersecret ./data_clef/masterseed.json \
    ${ARR[0]}


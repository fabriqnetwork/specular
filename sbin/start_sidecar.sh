#!/bin/sh
SBIN=`dirname $0`
SBIN="`cd "$SBIN"; pwd`"

# Check that the sidecar dotenv exists.
ENV=".sidecar.env"
if ! test -f $ENV; then
    echo "Expected dotenv at $ENV (does not exist)."
    exit
fi
echo "Using sidecar dotenv: $ENV"
. $ENV

if [ -z $SIDECAR_BIN ]; then
    # If no binary specified, assume repo directory structure.
    . $SBIN/configure.sh
fi

ARGS="
    --l1.endpoint $L1_ENDPOINT \
    --l2.endpoint $L2_ENDPOINT \
    --protocol.rollup-cfg-path $ROLLUP_CFG_PATH \
    --protocol.rollup-addr $ROLLUP_ADDR \
    --disseminator \
    --disseminator.addr \
    --disseminator.private-key $DISSEMINATOR_PRIVATE_KEY \
    --disseminator.sub-safety-margin $DISSEMINATOR_SUB_SAFETY_MARGIN \
    --disseminator.target-batch-size $DISSEMINATOR_TARGET_BATCH_SIZE \
    --validator \
    --validator.private-key $VALIDATOR_PRIVATE_KEY
"

echo "starting sidecar with the following flags:"
echo $ARGS

$SIDECAR $ARGS

#!/bin/bash
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

FLAGS=(
    "--l1.endpoint $L1_ENDPOINT"
    "--l2.endpoint $L2_ENDPOINT"
    "--protocol.rollup-cfg-path $ROLLUP_CFG_PATH"
    "--protocol.rollup-addr $ROLLUP_ADDR"
)

# Set disseminator flags.
if [ "$DISSEMINATOR" = true ] ; then
    echo "Enabling disseminator."
    DISSEMINATOR_PRIV_KEY=`cat "$DISSEMINATOR_PK_PATH"`
    FLAGS+=(
        "--disseminator"
        "--disseminator.private-key $DISSEMINATOR_PRIV_KEY"
        "--disseminator.sub-safety-margin $DISSEMINATOR_SUB_SAFETY_MARGIN"
        "--disseminator.target-batch-size $DISSEMINATOR_TARGET_BATCH_SIZE"
    )
fi
# Set validator flags.
if [ "$VALIDATOR" = true ] ; then
    echo "Enabling validator."
    VALIDATOR_PRIV_KEY=`cat "$VALIDATOR_PK_PATH"`
    FLAGS+=(
        "--validator"
        "--validator.private-key $VALIDATOR_PRIV_KEY"
    )
fi

echo "starting sidecar with the following flags:"
echo "${FLAGS[@]}"
$SIDECAR_BIN ${FLAGS[@]}

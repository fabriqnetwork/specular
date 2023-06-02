if [[ "$#" -ne 0 ]]; then
    echo "Usage: 
    USE_CLEF={true, false}
    IS_SEQUENCER={true, false}
    IS_VALIDATOR={true, false} $0"
    echo "Output: env variable SPECULAR_NODE_DEFAULTS"
    exit 1
fi

# SBIN=`dirname $0`
# SBIN="`cd "$SBIN"; pwd`"
# . $SBIN/configure.sh

# rollup config defaults
L1_ENDPOINT="ws://localhost:8545"
L1_CHAIN_ID=31337
L1_ROLLUP_GENESIS_BLOCK=0
L1_SEQUENCER_INBOX_ADDR="0x2E983A1Ba5e8b38AAAeC4B440B9dDcFBf72E15d1"
L1_ROLLUP_ADDR="0xF6168876932289D073567f347121A267095f3DD6"
L2_ENDPOINT="ws://0.0.0.0:4012"
L2_CLEF_ENDPOINT="http://127.0.0.1:8550"
VALIDATOR_STAKE_AMOUNT=100

SEQUENCER_ADDR="0xf39fd6e51aad88f6f4ce6ab8827279cfffb92266"
VALIDATOR_ADDR="0x70997970c51812dc3a010c7d01b50e0d17dc79c8"
# INDEXER_ADDR="0x3c44cdddb6a900fa2b585dd299e03d12fa4293bc"

# Start with common defaults
SPECULAR_NODE_DEFAULTS=(
    --rollup.l1.endpoint $L1_ENDPOINT
    --rollup.l1.chainid $L1_CHAIN_ID
    --rollup.l1.rollup-genesis-block $L1_ROLLUP_GENESIS_BLOCK
    --rollup.l1.sequencer-inbox-addr $L1_SEQUENCER_INBOX_ADDR
    --rollup.l1.rollup-addr $L1_ROLLUP_ADDR
    --rollup.l2.endpoint $L2_ENDPOINT
)

if $USE_CLEF ; then
    SPECULAR_NODE_DEFAULTS+=( --rollup.l2.clef-endpoint $CLEF_ENDPOINT )
fi

if $IS_SEQUENCER ; then
    SPECULAR_NODE_DEFAULTS+=( --rollup.sequencer.addr $SEQUENCER_ADDR )
fi

if $IS_VALIDATOR ; then 
    SPECULAR_NODE_DEFAULTS+=(
	--rollup.validator.addr $VALIDATOR_ADDR
	--rollup.validator.stake-amount $VALIDATOR_STAKE_AMOUNT
    )
fi

export SPECULAR_NODE_DEFAULTS
echo "Exported flags to SPECULAR_NODE_DEFAULTS"

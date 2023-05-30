SBIN_DIR=`dirname $0`
SBIN_DIR="`cd "$SBIN"; pwd`"

export NETWORK_ID=13527
export L1_ENDPOINT=ws://localhost:8545
export L1_CHAIN_ID=31337
export CLEF_ENDPOINT=http://127.0.0.1:8550

export SEQUENCER_INBOX_ADDR=0x2E983A1Ba5e8b38AAAeC4B440B9dDcFBf72E15d1
export ROLLUP_ADDR=0xF6168876932289D073567f347121A267095f3DD6

export SEQUENCER_ADDR=0xf39fd6e51aad88f6f4ce6ab8827279cfffb92266
export VALIDATOR_ADDR=0x70997970c51812dc3a010c7d01b50e0d17dc79c8
export INDEXER_ADDR=0x3c44cdddb6a900fa2b585dd299e03d12fa4293bc

export ROLLUP_STAKE_AMOUNT=100
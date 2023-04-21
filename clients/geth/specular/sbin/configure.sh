# Define directory structure for other scripts.
SBIN_DIR=`dirname $0`
SBIN_DIR="`cd "$SBIN"; pwd`"
export GETH_SPECULAR_DIR=$SBIN/../
export CLIENTS_DIR=$GETH_SPECULAR_DIR/../../
export ROOT_DIR=$CLIENTS_DIR/../
export CONTRACTS_DIR=$ROOT_DIR/contracts/
export DATA_DIR=$GETH_SPECULAR_DIR/data/

# Define network config
export NETWORK_ID=13527
export L1_ENDPOINT=ws://localhost:8545
export L1_CHAIN_ID=31337
export CLEF_ENDPOINT=http://127.0.0.1:8550

export COINBASE_ADDR=f39fd6e51aad88f6f4ce6ab8827279cfffb92266
export SEQUENCER_INBOX_ADDR=0xe7f1725E7734CE288F8367e1Bb143E90bb3F0512
export SEQUENCER_ADDR=0xf39fd6e51aad88f6f4ce6ab8827279cfffb92266
export ROLLUP_ADDR=0x5FC8d32690cc91D4c39d9d3abcBD16989F875707

export ROLLUP_STAKE_AMOUNT=100

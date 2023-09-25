# Define directory structure for other scripts.
SBIN=`dirname $0`
SBIN="`cd "$SBIN"; pwd`"
export GETH_SPECULAR_DIR=$SBIN/../
export CLIENTS_DIR=$GETH_SPECULAR_DIR/../../
export ROOT_DIR=$CLIENTS_DIR/../
export CONTRACTS_DIR=$ROOT_DIR/contracts/
export DATA_DIR=$GETH_SPECULAR_DIR/data/

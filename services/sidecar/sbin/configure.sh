# Define directory structure for other scripts.
SBIN=`dirname $0`
SBIN="`cd "$SBIN"; pwd`"
export ROOT_DIR=$SBIN/../
export CONTRACTS_DIR=$ROOT_DIR/contracts
export DATA_DIR=$ROOT_DIR/data
export GETH_DIR=$ROOT_DIR/services/el_client/go-ethereum
export SIDECAR_DIR=$ROOT_DIR/services/sidecar

export GETH_SPECULAR_DIR=$SBIN/../
export CLIENTS_DIR=$GETH_SPECULAR_DIR/../../

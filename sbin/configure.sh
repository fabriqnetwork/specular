# Define directory structure for other scripts.
SBIN=`dirname $0`
SBIN="`cd "$SBIN"; pwd`"
export ROOT_DIR=$SBIN/../
export CONTRACTS_DIR=$ROOT_DIR/contracts/
export CONFIG_DIR=$ROOT_DIR/config/
export SPECULAR_CLIENT_DIR=$ROOT_DIR/clients/geth/specular/ # to be moved out
export DATA_DIR=$SPECULAR_CLIENT_DIR/data/ # use config/data ?
export USE_CLEF=0

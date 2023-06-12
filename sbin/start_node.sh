#!/bin/bash
SBIN=`dirname $0`
SBIN="`cd "$SBIN"; pwd`"
. $SBIN/configure.sh

# Wrapping configuration scripts in this func removes positional args.
configure_args() {
    echo "Configuring..."
    . $SBIN/configure_geth_args.sh # sets GETH_ARGS
    . $SBIN/configure_specular_node_defaults.sh # sets SPECULAR_NODE_DEFAULT_ARGS
}

# Configure geth args
configure_args
echo "Starting specular node with args..." 
echo "[geth args]: ${GETH_ARGS[*]}"
echo "[specular args]: from the provided --rollup.config, if any (overriding defaults)"
cd $DATA_DIR && $SPECULAR_CLIENT_DIR/build/bin/geth "${GETH_ARGS[@]}" "${SPECULAR_NODE_DEFAULT_ARGS[@]}" "$@"

services/sidecar/build/bin/sidecar
services/cl_clients/magi/target/debug/magi
services/el_clients/go-ethereum/build/bin/geth

L1_GETH_BIN=$GETH_DIR/build/bin/geth # TODO: use l1, not sp?
SP_GETH_BIN=$GETH_DIR/build/bin/geth
SP_MAGI_BIN=$MAGI_DIR/target/debug/magi
SIDECAR_BIN=$SIDECAR_DIR/build/bin/sidecar



only inject, do not add to image

sequencer_pk.tx
validator_pk.txt
jwt_secret.txt
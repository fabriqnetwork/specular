#!/bin/bash
SBIN=`dirname $0`
SBIN="`cd "$SBIN"; pwd`"
. $SBIN/configure.sh
. $SBIN/configure_system.sh
cd $DATA_DIR

args=(
    --datadir ./data_indexer
    --http --http.addr '0.0.0.0' --http.port 4021 --http.api 'personal,eth,net,web3,txpool,miner,proof,debug'
    --ws --ws.addr '0.0.0.0' --ws.port 4022 --ws.api 'personal,eth,net,web3,txpool,miner,proof,debug'
    --http.corsdomain '*' --ws.origins '*'
    --http.vhosts '*'
    --networkid $NETWORK_ID
    --gcmode=archive
    --port 30305
    --authrpc.port 8562
    # Rollup flags
    --rollup.l1.endpoint $L1_ENDPOINT
    --rollup.l1.chainid $L1_CHAIN_ID
    --rollup.l1.rollup-genesis-block 0
    --rollup.l1.sequencer-inbox-addr $SEQUENCER_INBOX_ADDR
    --rollup.l1.rollup-addr $ROLLUP_ADDR
    # Validator with no active validation functionality enabled.
    # Driver flags
    # --rollup.driver.step-interval 2
    # --rollup.driver.retry-delay 8
    # --rollup.driver.num-attempts 4
)

if [ $USE_CLEF = true ]; then
    args+=(--rollup.l2.clef-endpoint $CLEF_ENDPOINT)
else
    args+=(--password ./password.txt)
fi

$GETH_SPECULAR_DIR/build/bin/geth "${args[@]}"

if [[ "$#" -ne 0 ]]; then
    echo "Usage: 
    DOCKER={0, 1}
    IS_SEQUENCER={0, 1}
    IS_VALIDATOR={0, 1}
    IS_INDEXER={0, 1} $0"
    echo "Output: env variable GETH_ARGS"
    exit 1
fi

DOCKER="${DOCKER:-0}"
IS_SEQUENCER="${IS_SEQUENCER:-0}"
IS_INDEXER="${IS_INDEXER:-0}"
IS_VALIDATOR="${IS_VALIDATOR:-0}"

echo "Building geth args..."

if [ $IS_SEQUENCER -eq "1" ] || [ $IS_VALIDATOR -eq "1" ] ; then
    API='personal,eth,net,web3,txpool,miner,proof,debug'
else
    API='eth,web3,txpool,debug'
fi

NETWORK_ID=13527
# Start with common args
GETH_ARGS=(
    --http
    --http.addr '0.0.0.0'
    --http.corsdomain '*'
    --http.api $API
    --ws
    --ws.addr '0.0.0.0'
    --ws.origins '*'
    --ws.api $API
    --networkid $NETWORK_ID
    --password "./password.txt"
)

if [ "$IS_INDEXER" -eq "1" ] ; then
    echo "Adding indexer args"
    GETH_ARGS+=(
	--http.vhosts='*'
	--gcmode=archive
    )
fi

if [ "$DOCKER" -eq "1" ] ; then
    echo "Adding docker args"
    GETH_ARGS+=(
	--http.port=8545
	--ws.port=8546
	--datadir "."
	--nodiscover
	--maxpeers 0
    )
else
    # Running locally so we use different ports and dirs
    if [ "$IS_SEQUENCER" -eq "1" ] ; then
	echo "Adding sequencer args"
	GETH_ARGS+=(
	    --datadir "./data_sequencer"
	    --http.port 4011
	    --ws.port 4012
	)
    elif [ "$IS_VALIDATOR" -eq "1" ] ; then
	echo "Adding validator args"
	GETH_ARGS+=(
	    --datadir "./data_validator"
	    --http.port 4018
	    --ws.port 4019
	    --port 30304
	    --authrpc.port 8561
	)
    elif [ "$IS_INDEXER" -eq "1" ] ; then
	echo "Adding indexer args"
	GETH_ARGS+=(
	    --datadir "./data_indexer"
	    --http.port 4021
	    --ws.port 4022
	    --port 30305
	    --authrpc.port 8562
	)
    else
	echo "No node type provided"
    fi
fi

echo "Assigned args to GETH_ARGS"

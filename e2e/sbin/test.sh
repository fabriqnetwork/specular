#!/bin/bash

# Configure variables
SBIN_DIR=`dirname $0`
SBIN_DIR="`cd "$SBIN_DIR"; pwd`"
set -o allexport
source $SBIN_DIR/configure.sh
set +o allexport

# Spin up L1 node
cd $CONTRACTS_DIR
docker run -d \
  --name geth_container \
  -v ./docker:/root \
  -p 8545:8545 \
  ethereum/client-go \
  --dev --dev.period 5 \
  --verbosity 5 \
  --http --http.api eth,web3,net --http.addr 0.0.0.0 \
  --ws --ws.api eth,net,web3 --ws.addr 0.0.0.0 --ws.port 8545

# declare -a Addresses( !
# '0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266' \
# '0x70997970C51812dc3A010C7d01b50e0d17dc79C8' \
# '0x3C44CdDdB6a900fa2b585dd299e03d12FA4293BC' \
# '0x90F79bf6EB2c4f870365E785982E1f101E93b906' \
# '0x15d34AAf54267DB7D7c367839AAf71A00a2C6A65' \
# '0x9965507D1a55bcC2695C58ba16FB37d819B0A4dc' \
# '0x976EA74026E726554dB657fA54763abd0C3a0aa9' \
# '0x14dC79964da2C08b23698B3D3cc7Ca32193d9955' \
# '0xa0Ee7A142d267C1f36714E4a8F75612F20a79720' \
# )

docker exec -it geth_container geth attach --exec \
  "eth.sendTransaction({ from: eth.coinbase, to: '0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266', value: web3.toWei(10, 'ether') })" \
  ws://host.docker.internal:8545

docker exec -it geth_container geth attach --exec \
  "eth.sendTransaction({ from: eth.coinbase, to: '0x70997970C51812dc3A010C7d01b50e0d17dc79C8', value: web3.toWei(10, 'ether') })" \
  ws://host.docker.internal:8545

docker exec -it geth_container geth attach --exec \
  "eth.sendTransaction({ from: eth.coinbase, to: '0x3C44CdDdB6a900fa2b585dd299e03d12FA4293BC', value: web3.toWei(10, 'ether') })" \
  ws://host.docker.internal:8545

docker exec -it geth_container geth attach --exec \
  "eth.sendTransaction({ from: eth.coinbase, to: '0x90F79bf6EB2c4f870365E785982E1f101E93b906', value: web3.toWei(10, 'ether') })" \
  ws://host.docker.internal:8545

npx hardhat deploy --network localhost

# Spin up L2 node
cd $PROJECT_DATA_DIR
$SBIN_DIR/sequencer.sh > $PROJECT_LOG_DIR/l2.log 2>&1 &
L2GETH_PID=$!

# Wait for nodes
$SBIN_DIR/wait-for-it.sh -t 60 $HOST:$L1_WS_PORT
$SBIN_DIR/wait-for-it.sh -t 60 $HOST:$L2_HTTP_PORT

# Run testing script
cd $CONTRACTS_DIR
npx hardhat deploy --network specularLocalDev
sleep 10

case $1 in

  general)
    npx ts-node scripts/testing.ts
    RESULT=$?
    ;;

  deposit)
    npx hardhat run scripts/bridge/standard_bridge_deposit_eth.ts
    RESULT=$?
    ;;

  withdraw)
    npx hardhat run scripts/bridge/standard_bridge_withdraw_eth.ts
    RESULT=$?
    ;;

  erc20)
    npx hardhat run scripts/bridge/standard_bridge_erc20.ts
    RESULT=$?
    ;;

  *)
    echo "unknown test"
    ;;
esac


# Kill nodes
disown $L2GETH_PID
disown $GANACHE_PID
kill $L2GETH_PID
kill $GANACHE_PID

# Clean up
$SBIN_DIR/clean.sh

exit $RESULT

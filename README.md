# Specular Monorepo

**Warning**: This repository is a prototype and should not be used in production yet.

## Directory Structure

<pre>
├── <a href="./services/">services</a>: L2 services
│   ├── <a href="./services/cl_clients">cl_clients</a>: Consensus-layer clients
│   ├── <a href="./services/el_clients/">el_clients</a>: Execution-layer clients
│   │      └── <a href="./services/el_clients/go-ethereum/">go-ethereum</a>: Minimally modified geth fork
│   └── <a href="./services/sidecar/">sidecar</a>: Sidecar services
├── <a href="./contracts">contracts</a>: L1 and L2 contracts
└── <a href="./lib/">lib</a>: Libraries used in L2 EL Clients
    └── <a href="./lib/el_golang_lib/">el_golang_lib</a>: Library for golang EL clients
</pre>

## License

Unless specified in subdirectories, this repository is licensed under the [Apache License 2.0](https://www.apache.org/licenses/LICENSE-2.0). See `LICENSE` for details.

## Build from source
```sh
git clone https://github.com/specularl2/specular
cd specular
git submodule update --init
```

Install all dependencies and build the node binaries.
Note: the rest of the commands below assume you are in the project root directory.
```
pnpm install
make
```

## Running a local devnet

This guide will walk you through how to set up a local devnet containing an L2 sequencer running over a local L1 node.

### Configure network

To configure a local devnet, you can just use an existing example from `config` as-is.
```sh
mkdir workspace
cp -a config/local_devnet/. workspace/ # copy all config files
```

This copies multiple dotenv files (below) which are expected by scripts to be in the current directory.
Some of these env files also reference the `genesis.json` and `rollup.json` used to configure the protocol.
```sh
.genesis.env   # Expected by `start_l1.sh`, `deploy_l1_contracts.sh` (not necessary for existing chains)
.contracts.env # Expected by `deploy_l1_contracts.sh`
.sp_geth.env   # Expected by `start_sp_geth.sh`
.sp_magi.env   # Expected by `start_sp_magi.sh`
.sidecar.env   # Expected by `start_sidecar.sh`
```

### Start L1
Run the below script to initialize a new L1 chain.
```sh
cd workspace
# Terminal #1
../sbin/start_l1.sh
```

### Start an L2 node
Deploy the L1 contracts on the newly started chain, and spin up all services required to run an L2 node.
```sh
# Terminal #2
../sbin/deploy_l1_contracts.sh && ../sbin/start_sp_geth.sh
# Terminal #3
../sbin/start_sp_magi.sh
# Terminal #4
../sbin/start_sidecar.sh
```

At this point, you'll have two chains started with the following parameters
- L2: chain ID `13527`, with a sequencer exposed on ports `4011` (http) and `4012` (ws).
- L1: chain ID `31337`, on port `8545` (ws).
To clear L2 network state, run `../sbin/clean_sp_geth.sh`.
To clear the L1 deployment, run `../sbin/clean_deployment.sh`.

### Transact using MetaMask

After the nodes are running, you can use your wallet (e.g. MetaMask) to send transactions to the sequencer, and see how transactions are executed, sequenced and confirmed.

**Configure wallet**

1. Go to `.sidecar.env` and copy the validator key to MetaMask. The account is pre-funded on L2, so you can use it to transact.
2. In `Settings -> Networks`, create a new network called `L2`, which connects to the sequencer.
Enter `http://localhost:4011` for the RPC URL, `13527` for the chain ID and `ETH` for currency symbol.

**Transact**

Now, you can use the pre-funded account to send transactions.
After an L2 transaction, in the Hardhat node console, observe the resulting L1 transactions:
- sequencer calls `appendTxBatch` to sequence transaction
- sequencer calls `createAssertion` to commit to a new disputable state assertion
- sequencer calls `confirmFirstUnresolvedAssertion` to confirm the assertion after every staker has attested to it.

If you restart the network after having transacted with MetaMask, don't forget to reset your MetaMask account.
Select the appropriate account, go to `Setting -> Advanced`, and click `Reset Account`.
This ensures the account nonce cache in MetaMask is cleared.

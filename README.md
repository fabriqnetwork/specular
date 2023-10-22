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

This guide will walk you through how to set up a local devnet containing an L2 sequencer running over a local L1 node.
In this example, all nodes operate honestly (no challenges are issued).

## Build from source
Clone this repository along with its submodules.
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

## Running a local network

### Configure network

The scripts in the next subsection expect to find the following dotenv files in the current working directory.
```sh
.sp_geth.env # Expected by `start_geth.sh`
.sp_magi.env # Expected by `start_sp_magi.sh`
.sidecar.env # Expected by `start_sidecar.sh` and `start_l1.sh`.
```

To configure a local devnet, you can just use the existing example in `local_devnet` as-is.
```sh
# Copy config files
cp config/deployments/local_devnet/* .
# Generate the genesis file
./sbin/create_genesis.sh
```

The env files copied set:
- L2 chain ID `13527`, with a sequencer exposed on ports `4011` (http) and `4012` (ws).
- L1 chain ID `31337`, on port `8545` (ws).

### Start a node

```sh
# Terminal #1: start L1 network
./sbin/start_l1.sh
# Terminal #2: start L2-EL client
./sbin/start_geth.sh
# Terminal #3: start L2-CL client
./sbin/start_magi.sh
# Terminal #4: start sidecar
./sbin/start_sidecar.sh
```

**Restart network**

To clear network state, run `./sbin/clean.sh`.

### Transact using MetaMask

After the nodes are running, you can use your wallet (e.g. MetaMask) to send transactions to the sequencer, and see how transactions are executed, sequenced and confirmed.

**Configure wallet**

1. Go to `$DATA_DIR/keys`, import the sequencer key to MetaMask.
Both accounts are pre-funded with 10 ETH each on L2 network, and you can use them to send transactions. Note: on L2, these two accounts are just normal accounts; not to be confused with the sequencer roles on L1 (the addresses are just being reused).
2. In `Settings -> Networks`, create a new network called `L2` which connects to the sequencer.
The sequencer node should be running while creating the network.
Enter `http://localhost:4011` for RPC URL, `13527` for Chain ID, `ETH` for currency symbol (we haven't changed the symbol yet).

**Transact**

Remember to reset the account after every clean start of the network.
Select the appropriate account, go to `Setting -> Advanced`, and click `Reset Account`.
This ensures the account nonce cache in MetaMask is cleared.

Now, you can use the pre-funded account to send transactions.

After an L2 transaction, in the Hardhat node console, observe the resultant transactions occuring on L1:
- sequencer calls `appendTxBatch` to sequence transaction
- sequencer calls `createAssertion` to create disputable assertion
- sequencer calls `confirmFirstUnresolvedAssertion` to confirm the assertion after every staker has attested to it.

If you restart the network after having transacted with MetaMask, don't forget to reset your MetaMask account.

## License

Unless specified in subdirectories, this repository is licensed under the [Apache License 2.0](https://www.apache.org/licenses/LICENSE-2.0). See `LICENSE` for details.

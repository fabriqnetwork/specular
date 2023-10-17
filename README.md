# Specular Monorepo

## Directory Structure

<pre>
├── <a href="./services/">services</a>: Specular L2 clients
│   ├── <a href="./services/cl_clients">cl_clients</a>: Consensus layer clients
│   ├── <a href="./services/el_clients/">el_clients</a>: Execution layer clients
│   │      └── <a href="./services/el_clients/go-ethereum/">go-ethereum</a>: Minimally modified geth fork
│   └── <a href="./services/sidecar/">sidecar</a>: The Specular sidecar service
├── <a href="./contracts">contracts</a>: Specular L1 and L2 contracts
└── <a href="./lib/">lib</a>: Libraries used in L2 EL Clients
    └── <a href="./lib/el_golang_lib/">el_golang_lib</a>: Library for golang clients
</pre>

## License

Unless specified in subdirectories, this repository is licensed under the [Apache License 2.0](https://www.apache.org/licenses/LICENSE-2.0). See `LICENSE` for details.

## Running a local network

This guide will demonstrate how to set up a rollup network containing L2 sequencer nodes, running over a Geth L1 node---all on your local machine.
After the 3 nodes are running, you can use MetaMask to send custom transactions to the sequencer, and see how transactions are executed on the L2 network, sequenced to the L1 network, and confirmed.
In this example, all nodes operate honestly (no challenges are issued).

### Build
Install all dependencies and build the modified L2 Geth node

```sh
cd SPECULAR_REPO
pnpm install
make geth sidecar
```

### Generate the genesis file

```sh
SPECULAR_REPO/sbin/create_genesis.sh
```

### L2 setup

```sh
SPECULAR_REPO/sbin/init_geth.sh
```

### L1 local dev node installation

See [here](https://github.com/SpecularL2/specular/tree/main/contracts) for more details.

### Start nodes

```sh
# Terminal #1: start L1 node
SPECULAR_REPO/sbin/start_l1.sh

# Terminal #2: start L2 EL client
SPECULAR_REPO/sbin/start_geth.sh

# Terminal #3: start Specular sidecar
SPECULAR_REPO/sbin/start_sidecar.sh
```

Make sure there are logs for `Sequencer started` in the respective console.
In the first terminal where L1 node is running, you can see the sequencer staked on the Rollup contract.

**Restarts**

Currently, the sequencer must start in a clean environment; i.e. you need to clean and reinitialize both L1 and L2 on every start.

To restart the L1 node, use `Ctrl-C` to stop the current running one and run `npx hardhat node` again.

To reinitialize L2 node, under `sbin` directory, run `./clean.sh && ./init.sh`.

Do not forget to reset MetaMask account if you have sent some transactions on L2 (see below for more details).

### Transact using MetaMask

**Configuration**

1. Go to `data/keys`, import the sequencer key to MetaMask.
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
- sequencer calls `confirmFirstUnresolvedAssertion` to confirm the assertion after every stakers staked on the assertion

*Make sure that sequencer node is started before sending any transaction to L2.*

### Scenario Parameters

L1: Hardhat, chain ID `31337`, http/ws on port `8545`.

L2: Chain ID `13527`. Sequencer: http on port `4011`, ws on port `4012`; Validator: http on port `4018`, ws on port `4019`.

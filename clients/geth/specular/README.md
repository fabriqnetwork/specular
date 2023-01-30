# Specular Client

## Build

First setup L1 and then follow following steps
```sh
make install
```

## Running a local network

This guide will demonstrate how to set up a rollup network containing L2 sequencer validator nodes, running over a Hardhat L1 node---all on your local machine.
After the 3 nodes are running, you can use MetaMask to send custom transactions to the sequencer, and see how transactions are executed on the L2 network, sequenced to the L1 network, and validated and confirmed.
In this example, all nodes operate honestly (no challenges are issued).

### L2 setup

```sh
cd sbin
./import_accounts.sh
./init.sh
```

### L1 local dev node installation

See [here](https://github.com/SpecularL2/specular/tree/main/contracts) for more details.

### Start nodes

```sh
# Terminal #1: start L1 node
cd contracts
yarn install
npx hardhat node

# Terminal #2: start sequencer
cd clients/geth/specular/sbin
./start_sequencer.sh

# Terminal #3: start validator
cd clients/geth/specular/sbin
./start_validator.sh
```

Make sure there are logs for `Sequencer started` and `Validator started` in the respective consoles.
In the first terminal where L1 node is running, you can see both sequencer and validator are staked on the Rollup contract.

**Restarts**

Currently, the sequencer and validator must start in a clean environment; i.e. you need to clean and reinitialize both L1 and L2 on every start.

To restart the L1 node, use `Ctrl-C` to stop the current running one and run `npx hardhat node` again.

To reinitialize L2 node, under `rollup/test-node` directory, run `./clean.sh && ./init.sh`.

Do not forget to reset MetaMask account if you have sent some transactions on L2 (see below for more details).

### Transact using MetaMask

**Configuration**

1. Go to `test-node/keys`, import the sequencer and validator keys to MetaMask.
Both accounts are pre-funded with 10 ETH each on L2 network, and you can use them to send transactions. Note: on L2, these two accounts are just normal accounts; not to be confused with the sequencer/validator roles on L1 (the addresses are just being reused).
2. In `Settings -> Networks`, create a new network called `L2` which connects to the sequencer.
The sequencer node should be running while creating the network.
Enter `http://localhost:4011` for RPC URL, `13527` for Chain ID, `ETH` for currency symbol (we haven't changed the symbol yet).

**Transact**

Remember to reset the account after every clean start of the network.
Select the appropriate account, go to `Setting -> Advanced`, and click `Reset Account`.
This ensures the account nonce cache in MetaMask is cleared.

Now, you can use the pre-funded account to send transactions.
For example, you can send 1 ETH from the sequencer account to the validator account on L2.

After an L2 transaction, in the Hardhat node console, observe the resultant transactions occuring on L1:
- sequencer calls `appendTxBatch` to sequence transaction
- sequencer calls `createAssertion` to create disputable assertion
- validator calls `advanceStake` after validating the assertion
- sequencer calls `confirmFirstUnresolvedAssertion` to confirm the assertion after every stakers staked on the assertion

*Make sure that sequencer node and validator node are all started before sending any transaction to L2.*

### Scenario Parameters

L1: Hardhat, chain ID `31337`, http/ws on port `8545`.

L2: Chain ID `13527`. Sequencer: http on port `4011`, ws on port `4012`; Validator: http on port `4018`, ws on port `4019`.

## License

The Golang binding of Specular L2 contracts (files under `bindings` directory) is licensed under the [Apache License 2.0](https://www.apache.org/licenses/LICENSE-2.0), see `LICENSE` for details.

The Specular client binary (files under `cmd` directory) is directly modified from of the [go-ethereum](https://github.com/ethereum/go-ethereum) binary. It is licensed under the [GNU General Public License v3.0](https://www.gnu.org/licenses/gpl-3.0.en.html) in accordance with the [original license](https://github.com/ethereum/go-ethereum/blob/master/COPYING), see `COPYING` for details. Major modifications are marked with `<specular modifications>`.

Files under `internal` directory are directly copied from the [go-ethereum](https://github.com/ethereum/go-ethereum) library, which is [originally licensed](https://github.com/ethereum/go-ethereum/blob/master/COPYING.LESSER) under the [GNU Lesser General Public License v3.0](https://www.gnu.org/licenses/lgpl-3.0.en.html), see `COPYING` for details.

The Specular proof module (files under `proof` directory) is licensed under the [Apache License 2.0](https://www.apache.org/licenses/LICENSE-2.0), see `LICENSE` for details.

The Specular rollup module (files under `rollup` directory) is licensed under the [GNU Lesser General Public License v3.0](https://www.gnu.org/licenses/lgpl-3.0.en.html), see `COPYING.LESSER` for details.

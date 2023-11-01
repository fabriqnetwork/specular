# Specular

Specular is an L2 system designed to scale Ethereum securely, with minimal additional trust assumptions. Specifically, it is an EVM-native optimistic rollup, relying on existing Ethereum infrastructure both to bootstrap protocol security and to enable native compatibility for all existing Ethereum applications & tooling.

This repository contains the L1 protocol contracts and L2 node software. The source is licensed under the [Apache License 2.0](https://www.apache.org/licenses/LICENSE-2.0) (unless otherwise specified)—see <a href="./LICENSE.md">`LICENSE`</a> for details. **Warning**: This repository is a prototype and should not be used in production yet.

## For developers

See <a href="./docs/system.md">system.md</a> for a system overview.

### Build from source

Install the following dependencies:
- [`Node.js v16`](https://nodejs.org/en/)
- [`pnpm`](https://pnpm.io/installation#using-corepack)
- [`foundry`](https://book.getfoundry.sh/getting-started/installation)

Then to build, run:
```sh
# Fetch the repository and its submodules.
git clone https://github.com/specularl2/specular
cd specular
git submodule update --init --recursive
# Install dependencies and build binaries
pnpm install
make
```

### Running a local devnet

This guide will walk you through how to set up a local devnet containing an L2 sequencer running over a local L1 network.
Note: the commands that follow below assume you are in the project root directory.

### Configure network

To configure a local devnet, you can just use an existing example from <a href="./config/">config`</a> as-is.
```sh
mkdir workspace
cp -a config/local_devnet/. workspace/ # copy all config files
```

This copies multiple dotenv files (below) which are expected by scripts to be in the current directory.
Some of these env files also reference the `genesis.json` and `rollup.json` used to configure the protocol.
```sh
.genesis.env   # Expected by `start_l1.sh` & `deploy_l1_contracts.sh` (not necessary for existing chains)
.contracts.env # Expected by `deploy_l1_contracts.sh`
.sp_geth.env   # Expected by `start_sp_geth.sh`
.sp_magi.env   # Expected by `start_sp_magi.sh`
.sidecar.env   # Expected by `start_sidecar.sh`
```

### Start L1
Run the below script to initialize a new local L1 chain.
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

## For users
To learn more, see [`specular.network`](https://specular.network/).

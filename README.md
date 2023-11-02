# Specular

Specular is an L2 system designed to scale Ethereum securely, with minimal additional trust assumptions. Specifically, it is an EVM-native optimistic rollup, relying on existing Ethereum infrastructure both to bootstrap protocol security and to enable native compatibility for all existing Ethereum applications & tooling.

This repository contains the L1 protocol contracts and L2 node software. The source is licensed under the [Apache License 2.0](https://www.apache.org/licenses/LICENSE-2.0) (unless otherwise specified)—see <a href="./LICENSE.md">LICENSE</a> for details.

**Warning**: This repository is a prototype and should not be used in production yet.

## For developers

See <a href="./docs/system.md">system.md</a> for a system overview and <a href="./docs/development.md">development.md</a> for further information on how to deploy the system, test changes and contribute to the repository.

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
pnpm install && make
```

### Running a local devnet

This section will walk you through how to set up a local devnet containing an L2 sequencer running over a local L1 network.
Note: the commands that follow below assume you are in the project root directory.

### Configure network

To configure a local devnet, you can just use an existing example from <a href="./config/">config</a> as-is.
```sh
mkdir workspace
cp -a config/local_devnet/. workspace/ # copy all config files
```

### Start L1
Run the below script to initialize a new local L1 chain.
```sh
cd workspace
../sbin/start_l1.sh # Terminal #1
```

### Start an L2 node
Deploy the L1 contracts on the newly started chain, and spin up all services required to run an L2 node.
```sh
../sbin/deploy_l1_contracts.sh && ../sbin/start_sp_geth.sh # Terminal #2
../sbin/start_sp_magi.sh # Terminal #3
../sbin/start_sidecar.sh # Terminal #4
```

At this point, you'll have two chains started with the following parameters
- L2: chain ID `13527`, with a sequencer exposed on ports `4011` (http) and `4012` (ws).
- L1: chain ID `31337`, on port `8545` (ws).
To re-run the network from clean state, run `../sbin/deploy_l1_contracts.sh -c && ../sbin/start_sp_geth.sh -c`.

## For users
To learn more, see [`specular.network`](https://specular.network/).

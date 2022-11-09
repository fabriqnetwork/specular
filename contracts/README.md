# Specular Contracts

This is the L1 contract of Specular. It is configured for both Hardhat and Founrdry.
Currently, Hardhat is used for deployment and Froundry is used for testing.

## Setup

Install following tools:

- [`Node.js v16`](https://nodejs.org/en/)
- [`yarn`](https://yarnpkg.com/getting-started/install)
- [`foundry`](https://book.getfoundry.sh/getting-started/installation)

Clone the repository:

```sh
git clone https://github.com/SpecularL2/specular
git submodule update --init
cd contracts
```

Install `npm` packages and git hooks:

```sh
yarn install
```

## Run Tests

```sh
npm run tests
```

## Run Lints

```sh
npm run lint:fix
```

## Local Slither Check

Install [`slither`](https://github.com/crytic/slither):

```sh
pip3 install slither-analyzer
```

Run slither:

```sh
slither .
```

## Generate Golang Bindings

```sh
./scripts/build_bindings.sh
```

## Run Local Development Node

```sh
npx hardhat node
```

Above command will start a Ethereum node serving as L1.
It can be accessed via `http://localhost:8545` or `ws://localhost:8545`.

It is configured to mine immediately when there is any transaction, or after 5 seconds idle.

As a convention, the first funded account is the sequencer, the second is the validator.

See `hardhat.config.js` for detailed configuration.

See `deploy/deploy.js` for how contracts are deployed and initialized.

## Attach to Local Development Node

```sh
npx hardhat console
```

Above command will start a Node.js console.

You can run `const provider = waffle.provider` in the console to obtain the ethers.js provider connected to the local development node.
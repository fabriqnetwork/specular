# OneStepProof End-to-End Test

Guides of a simple test to check the end-to-end functionality of OneStepProof.

## Local Environment Setup

Setup the local environment as described in the [L1 README](./README.md) and [L2 README](../clients/geth/specular/README.md).

## Start the L1 node

Start the L1 node and deploy the L1 contracts:
```sh
cd contracts
npx hardhat node
```

## Start the L2 node

In a separate terminal, start the L2 node:
```sh
cd clients/geth/specular/rollup/test-node
./import_accounts.sh # skip if you have already imported accounts
./clean.sh
./init.sh
./start_sequencer.sh
```

## Deploy the VerifierTestDriver contract

In a separate terminal, deploy the VerifierTestDriver contract:
```sh
cd contracts
npx hardhat --network localhost deployVerifierDriver
```

Record the address of the deployed `VerifierTestDriver` contract.

## Run the test

An example proof file is located at `contracts/test/osp/sample_proof.json`.

```sh
cd contracts
npx hardhat --network localhost verifyOsp --proof test/osp/sample_proof.json --addr <VERIFIER_TEST_DRIVER_ADDRESS>
```

You can see the `nextHash` logged in the *L1 node terminal* and compare it with the `nextHash` in the proof file.

## Generate Proofs

To generate proofs, you need to create some transactions on L2.
Please use remix to deploy some sample contracts and call some functions on them.
Currently proof generation through RPC calls does not support transfer transactions.

After you have created some transactions, you can generate proofs by running the following command:
```sh
cd contracts
npx hardhat generateOsp --hash <TRANSCATION_HASH> --step <STEP_NUMBER> --file <PROOF_FILE_TO_SAVE>
```

`<STEP_NUMBER>` must be greater or equal to 1.

**TODO**: This document is under-construction.

## Contribute
In addition to the dependencies listed in `README.md`, you'll also need to install:
- [`slither`](https://github.com/crytic/slither)
- [`Husky`](https://www.npmjs.com/package/husky)
- ... **TODO**

**TODO**

## Contract development
**TODO**
Protocol contracts (L1 contracts and L2 pre-deploys) are located under `contracts/`.
The package is configured to use both Hardhat (for deployment) and Founrdry (for testing).
See `hardhat.config.js` for the full configuration and `deploy/deploy.js` for how contracts are deployed.

Under `contracts/`, to run tests:
```sh
forge test
```
And to run slither:
```sh
slither .
```

## Testing

**TODO**

### Test system locally

**TODO**

To use the local deployment scripts in `sbin`, you'll need the following dotenv files.
```sh
.genesis.env   # Expected by `start_l1.sh` & `deploy_l1_contracts.sh` (not necessary for existing chains)
.contracts.env # Expected by `deploy_l1_contracts.sh`
.sp_geth.env   # Expected by `start_sp_geth.sh`
.sp_magi.env   # Expected by `start_sp_magi.sh`
.sidecar.env   # Expected by `start_sidecar.sh`
```
See `config/example/` for documentation for each dotenv.
Some of these env files also reference the `genesis.json` and `rollup.json` used to configure the protocol.

See `e2e.md` for running E2E tests.

### Test system manually
Given a running local devnet, you can transact with your wallet (e.g. MetaMask) to send transactions to the sequencer.

**Configure wallet**
1. Go to `.sidecar.env` and copy the validator key to MetaMask. The account is pre-funded on L2, so you can use it to transact.
2. In `Settings -> Networks`, create a new network called `L2`, which connects to the sequencer.
Enter `http://localhost:4011` for the RPC URL, `13527` for the chain ID and `ETH` for currency symbol.

**Transact**
You can now use the pre-funded account to send transactions.
After an L2 transaction, in the Hardhat node console, observe the resulting L1 transactions:
- sequencer calls `appendTxBatch` to sequence transaction
- sequencer calls `createAssertion` to commit to a new disputable state assertion
- sequencer calls `confirmFirstUnresolvedAssertion` to confirm the assertion after every staker has attested to it.

If you restart the network after having transacted with MetaMask, don't forget to reset your MetaMask account.
Select the appropriate account, go to `Setting -> Advanced`, and click `Reset Account`.
This ensures the account nonce cache in MetaMask is cleared.

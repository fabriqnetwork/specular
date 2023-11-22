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
.contracts.env # Expected by `deploy_l1_contracts.sh` (not necessary for existing chains)
.sp_geth.env   # Expected by `start_sp_geth.sh`
.sp_magi.env   # Expected by `start_sp_magi.sh`
.sidecar.env   # Expected by `start_sidecar.sh`
```
See `config/example/` for documentation for each dotenv.
Some of these env files also reference the `genesis.json` and `rollup.json` used to configure the protocol.

See `e2e.md` for running E2E tests.

### Test system manually
Given a running local devnet, you can transact with your wallet (e.g. using `cast` or MetaMask) to send transactions to the sequencer. For example:
```bash
cast send \
    --rpc-url http://localhost:4011 \
    --chain 13527 \
    --private-key `cat validator_pk.txt` \
    --value 0.01ether \
    0x0000000000000000000000000000000000000000 # to-address
```

After an L2 transaction, in the L1 node console, observe the resulting L1 transactions:
- disseminator calls `appendTxBatch` to sequence transaction
- validator calls `createAssertion` to commit to a new disputable state assertion
- validator calls `confirmFirstUnresolvedAssertion` to confirm the assertion after every staker has attested to it.

### Troubleshooting

If you see a message like`Forkchoice requested unknown head` logged by `sp-geth`, it may be because it's using stale data, while the CL client is using the correct genesis hash.
- This can happen due to re-running geth without cleaning. You can use `start_sp_geth.sh -c` to do so.
- This may also happen if you called `start_sp_geth.sh` without waiting for `start_l1.sh` to finish creating the new deployment configs. Make sure you wait a couple seconds—you should see `L1 started... (Use ctrl-c to stop)`.

# Protocol contracts
Protocol contracts are located in `contracts/`. 
These are the L1 contracts and L2 pre-deploys of Specular. 

The package is configured to use both Hardhat (for deployment) and Founrdry (for testing).
See `hardhat.config.js` for the full configuration and `deploy/deploy.js` for how contracts are deployed.

### Run tests

```sh
forge test
```

### Local slither check
Install [`slither`](https://github.com/crytic/slither):
```sh
pip3 install slither-analyzer
```

Run slither:
```sh
slither .
```

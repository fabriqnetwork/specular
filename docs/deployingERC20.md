# Deploying ERC20 Tokens

This is a quick guide on how to deploy L2 versions of ERC20 tokens already existing on L1.

## Deploying the L2 contract

The L2 token must implement the [`IMintableERC20`](../contracts/src/bridge/mintable/IMintableERC20.sol) interface.

The easiest way deploying the L2 token contract is by calling the `MintableERC20Factory` contract.
This is a predeploy at `0x2A000000000000000000000000000000000000F0`

See [this script](../contracts/scripts/bridge/deployERC20Token.ts) `MintableERC20Factory` for how to interact with the contract.
It can also be used as a CLI tool as follows:

```
npx ts-node ./scripts/bridge/deployERC20Token.ts --rpc http://127.0.0.1:4011 --name test --symbol TT --address 0x000000000000000000000000000000000000dead
```

## Adding the contract pair to the token list

Tokens can be added to the [Specular Token List](https://github.com/SpecularL2/specular-token-list) to be accessible through the Specular Bridge UI.

This is also how we keep track of canonical token mappings.

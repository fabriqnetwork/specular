# Deploying ERC20 Tokens

This is a quick guide on how to deploy L2 versions of ERC20 tokens already existing on L1.

## Deploying the L2 contract

The L2 token must implement the [`IMintableERC20`](../contracts/src/bridge/mintable/IMintableERC20.sol) interface.

The easiest way deploying the L2 contract is by calling the `MintableERC20Factory` contract.

TODO: add predeploy address here

See [this](../contracts/scripts/bridge/deployERC20Token.ts) `MintableERC20Factory`.

## Adding the contract pair to the token list

Tokens can be added to the [Specular Token List](https://github.com/SpecularL2/specular-token-list) to be accessible through the Specular Bridge UI.

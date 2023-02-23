# RollupLib
[Git Source](https://github.com/SpecularL2/specular/blob/559c78f8b09496c7f5c8f6e0b0262bee5e41d9a4/src/RollupLib.sol)


## Functions
### stateHash

Computes the hash of `execState`.


```solidity
function stateHash(ExecutionState memory execState) internal pure returns (bytes32);
```

## Structs
### ExecutionState

```solidity
struct ExecutionState {
    uint256 l2GasUsed;
    bytes32 vmHash;
}
```


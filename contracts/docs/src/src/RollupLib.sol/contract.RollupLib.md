# RollupLib
[Git Source](https://github.com/SpecularL2/specular/blob/c54213cfb14aca9d44e839341f672dd978834f68/src/RollupLib.sol)


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


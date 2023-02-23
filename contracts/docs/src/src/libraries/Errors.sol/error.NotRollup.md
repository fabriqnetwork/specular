# NotRollup
[Git Source](https://github.com/SpecularL2/specular/blob/559c78f8b09496c7f5c8f6e0b0262bee5e41d9a4/src/libraries/Errors.sol)

*Thrown when unauthorized (!rollup) address calls an only-rollup function*


```solidity
error NotRollup(address sender, address rollup);
```


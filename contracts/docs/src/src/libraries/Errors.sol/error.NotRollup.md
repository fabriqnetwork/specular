# NotRollup
[Git Source](https://github.com/SpecularL2/specular/blob/c54213cfb14aca9d44e839341f672dd978834f68/src/libraries/Errors.sol)

*Thrown when unauthorized (!rollup) address calls an only-rollup function*


```solidity
error NotRollup(address sender, address rollup);
```


# NotSequencer
[Git Source](https://github.com/SpecularL2/specular/blob/c54213cfb14aca9d44e839341f672dd978834f68/src/libraries/Errors.sol)

*Thrown when unauthorized (!sequencer) address calls an only-sequencer function*


```solidity
error NotSequencer(address sender, address sequencer);
```


# NotSequencer
[Git Source](https://github.com/SpecularL2/specular/blob/559c78f8b09496c7f5c8f6e0b0262bee5e41d9a4/src/libraries/Errors.sol)

*Thrown when unauthorized (!sequencer) address calls an only-sequencer function*


```solidity
error NotSequencer(address sender, address sequencer);
```


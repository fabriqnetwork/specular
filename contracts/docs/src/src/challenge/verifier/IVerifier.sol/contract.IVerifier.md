# IVerifier
[Git Source](https://github.com/SpecularL2/specular/blob/559c78f8b09496c7f5c8f6e0b0262bee5e41d9a4/src/challenge/verifier/IVerifier.sol)


## Functions
### verifyOneStepProof


```solidity
function verifyOneStepProof(IVerificationContext ctx, bytes32 currStateHash, bytes calldata encodedProof)
    external
    view
    returns (bytes32 nextStateHash);
```


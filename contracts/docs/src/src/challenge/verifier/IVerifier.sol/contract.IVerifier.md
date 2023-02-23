# IVerifier
[Git Source](https://github.com/SpecularL2/specular/blob/c54213cfb14aca9d44e839341f672dd978834f68/src/challenge/verifier/IVerifier.sol)


## Functions
### verifyOneStepProof


```solidity
function verifyOneStepProof(IVerificationContext ctx, bytes32 currStateHash, bytes calldata encodedProof)
    external
    view
    returns (bytes32 nextStateHash);
```


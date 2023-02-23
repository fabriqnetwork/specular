# Verifier
[Git Source](https://github.com/SpecularL2/specular/blob/559c78f8b09496c7f5c8f6e0b0262bee5e41d9a4/src/challenge/verifier/Verifier.sol)

**Inherits:**
[IVerifier](/src/challenge/verifier/IVerifier.sol/contract.IVerifier.md), Initializable, UUPSUpgradeable, OwnableUpgradeable


## Functions
### initialize


```solidity
function initialize() public initializer;
```

### constructor


```solidity
constructor();
```

### _authorizeUpgrade


```solidity
function _authorizeUpgrade(address) internal override onlyOwner;
```

### verifyOneStepProof


```solidity
function verifyOneStepProof(IVerificationContext, bytes32, bytes calldata)
    external
    pure
    override
    returns (bytes32 nextStateHash);
```


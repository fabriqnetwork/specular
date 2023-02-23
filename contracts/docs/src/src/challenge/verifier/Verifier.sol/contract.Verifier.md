# Verifier
[Git Source](https://github.com/SpecularL2/specular/blob/c54213cfb14aca9d44e839341f672dd978834f68/src/challenge/verifier/Verifier.sol)

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


# AssertionMap
[Git Source](https://github.com/SpecularL2/specular/blob/559c78f8b09496c7f5c8f6e0b0262bee5e41d9a4/src/AssertionMap.sol)


## State Variables
### assertions

```solidity
mapping(uint256 => Assertion) public assertions;
```


### rollupAddress

```solidity
address public rollupAddress;
```


## Functions
### rollupOnly


```solidity
modifier rollupOnly();
```

### constructor


```solidity
constructor(address _rollupAddress);
```

### getStateHash


```solidity
function getStateHash(uint256 assertionID) external view returns (bytes32);
```

### getInboxSize


```solidity
function getInboxSize(uint256 assertionID) external view returns (uint256);
```

### getParentID


```solidity
function getParentID(uint256 assertionID) external view returns (uint256);
```

### getDeadline


```solidity
function getDeadline(uint256 assertionID) external view returns (uint256);
```

### getProposalTime


```solidity
function getProposalTime(uint256 assertionID) external view returns (uint256);
```

### getNumStakers


```solidity
function getNumStakers(uint256 assertionID) external view returns (uint256);
```

### isStaker


```solidity
function isStaker(uint256 assertionID, address stakerAddress) external view returns (bool);
```

### createAssertion


```solidity
function createAssertion(uint256 assertionID, bytes32 stateHash, uint256 inboxSize, uint256 parentID, uint256 deadline)
    external
    rollupOnly;
```

### stakeOnAssertion


```solidity
function stakeOnAssertion(uint256 assertionID, address stakerAddress) external rollupOnly;
```

### deleteAssertion


```solidity
function deleteAssertion(uint256 assertionID) external rollupOnly;
```

## Errors
### ChildInboxSizeMismatch

```solidity
error ChildInboxSizeMismatch();
```

### SiblingStateHashExists

```solidity
error SiblingStateHashExists();
```

## Structs
### Assertion

```solidity
struct Assertion {
    bytes32 stateHash;
    uint256 inboxSize;
    uint256 parent;
    uint256 deadline;
    uint256 proposalTime;
    uint256 numStakers;
    mapping(address => bool) stakers;
    uint256 childInboxSize;
    mapping(bytes32 => bool) childStateHashes;
}
```


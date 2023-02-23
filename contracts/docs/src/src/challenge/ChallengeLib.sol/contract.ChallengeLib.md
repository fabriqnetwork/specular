# ChallengeLib
[Git Source](https://github.com/SpecularL2/specular/blob/c54213cfb14aca9d44e839341f672dd978834f68/src/challenge/ChallengeLib.sol)


## Functions
### initialBisectionHash

Computes the initial bisection hash.


```solidity
function initialBisectionHash(bytes32 startStateHash, bytes32 endStateHash, uint256 numSteps)
    internal
    pure
    returns (bytes32);
```
**Parameters**

|Name|Type|Description|
|----|----|-----------|
|`startStateHash`|`bytes32`|Hash of agreed-upon start state.|
|`endStateHash`|`bytes32`|Disagreed-upon end state.|
|`numSteps`|`uint256`|Number of steps from the end of `startState` to the end of `endState`.|


### computeBisectionHash

Computes H(bisection || segmentStart || segmentLength)


```solidity
function computeBisectionHash(
    bytes32[] memory bisection,
    uint256 challengedSegmentStart,
    uint256 challengedSegmentLength
) internal pure returns (bytes32);
```
**Parameters**

|Name|Type|Description|
|----|----|-----------|
|`bisection`|`bytes32[]`|Array of stateHashes. First element is the last agreed upon state hash.|
|`challengedSegmentStart`|`uint256`|The number of steps preceding `bisection[1]`, relative to the assertion being challenged.|
|`challengedSegmentLength`|`uint256`|Length of bisected segment (in steps), from the start of bisection[1] to the end of bisection[-1].|


### firstSegmentLength

Returns length of first segment in a bisection.


```solidity
function firstSegmentLength(uint256 length, uint256 bisectionDegree) internal pure returns (uint256);
```

### otherSegmentLength

Returns length of a segment (after first) in a bisection.


```solidity
function otherSegmentLength(uint256 length, uint256 bisectionDegree) internal pure returns (uint256);
```


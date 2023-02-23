# IChallenge
[Git Source](https://github.com/SpecularL2/specular/blob/559c78f8b09496c7f5c8f6e0b0262bee5e41d9a4/src/challenge/IChallenge.sol)

Protocol execution:
`initialize` (challenger, via Rollup) ->
`initializeChallengeLength` (defender) ->
`bisectExecution` (challenger, defender -- alternating) ->
`verifyOneStepProof`


## Functions
### initialize

Initializes contract.


```solidity
function initialize(
    address _defender,
    address _challenger,
    IVerifier _verifier,
    address _resultReceiver,
    bytes32 _startStateHash,
    bytes32 _endStateHash
) external;
```
**Parameters**

|Name|Type|Description|
|----|----|-----------|
|`_defender`|`address`|Defending party.|
|`_challenger`|`address`|Challenging party. Challenger starts.|
|`_verifier`|`IVerifier`|Address of the verifier contract.|
|`_resultReceiver`|`address`|Address of contract that will receive the outcome (via callback `completeChallenge`).|
|`_startStateHash`|`bytes32`|Bisection root being challenged.|
|`_endStateHash`|`bytes32`|Bisection root being challenged.|


### initializeChallengeLength

Initializes the length of the challenge. Must be called by defender before bisection rounds begin.


```solidity
function initializeChallengeLength(uint256 _numSteps) external;
```
**Parameters**

|Name|Type|Description|
|----|----|-----------|
|`_numSteps`|`uint256`|Number of steps executed from the start of the assertion to its end. If this parameter is incorrect, the defender will be slashed (assuming successful execution of the protocol by the challenger).|


### bisectExecution

Bisects a segment. The challenged segment is defined by: {`challengedSegmentStart`, `challengedSegmentLength`, `bisection[0]`, `oldEndHash`}


```solidity
function bisectExecution(
    bytes32[] calldata bisection,
    uint256 challengedSegmentIndex,
    bytes32[] calldata prevBisection,
    uint256 prevChallengedSegmentStart,
    uint256 prevChallengedSegmentLength
) external;
```
**Parameters**

|Name|Type|Description|
|----|----|-----------|
|`bisection`|`bytes32[]`|Bisection of challenged segment. Each element is a state hash (see `ChallengeLib.stateHash`). The first element is the last agreed upon state hash. Must be of length MAX_BISECTION_LENGTH for all rounds except the last. In the last round, the bisection segments must be single steps.|
|`challengedSegmentIndex`|`uint256`|Index into `prevBisection`. Must be greater than 0 (since the first is agreed upon).|
|`prevBisection`|`bytes32[]`|Bisection in the preceding round.|
|`prevChallengedSegmentStart`|`uint256`|Offset of the segment challenged in the preceding round (in steps). Note: this is relative to the assertion being challenged (i.e. always between 0 and the initial `numSteps`).|
|`prevChallengedSegmentLength`|`uint256`|Length of the segment challenged in the preceding round (in steps).|


### verifyOneStepProof

Verifies one step proof and completes challenge protocol.


```solidity
function verifyOneStepProof(
    bytes memory proof,
    uint256 challengedStepIndex,
    bytes32[] calldata prevBisection,
    uint256 prevChallengedSegmentStart,
    uint256 prevChallengedSegmentLength
) external;
```
**Parameters**

|Name|Type|Description|
|----|----|-----------|
|`proof`|`bytes`|TODO.|
|`challengedStepIndex`|`uint256`|Index into `prevBisection`. Must be greater than 0 (since the first is agreed upon).|
|`prevBisection`|`bytes32[]`|Bisection in the preceding round. Each segment must now be of length 1 (i.e. a single step).|
|`prevChallengedSegmentStart`|`uint256`|Offset of the segment challenged in the preceding round (in steps). Note: this is relative to the assertion being challenged (i.e. always between 0 and the initial `numSteps`).|
|`prevChallengedSegmentLength`|`uint256`|Length of the segment challenged in the preceding round (in steps).|


### timeout

Triggers completion of challenge protocol if a responder timed out.


```solidity
function timeout() external;
```

### currentResponder


```solidity
function currentResponder() external view returns (address);
```

### currentResponderTimeLeft


```solidity
function currentResponderTimeLeft() external view returns (uint256);
```

## Events
### ChallengeCompleted

```solidity
event ChallengeCompleted(address winner, address loser, CompletionReason reason);
```

### Bisected

```solidity
event Bisected(bytes32 challengeState, uint256 challengedSegmentStart, uint256 challengedSegmentLength);
```

## Enums
### CompletionReason

```solidity
enum CompletionReason {
    OSP_VERIFIED,
    TIMEOUT
}
```


# IRollup
[Git Source](https://github.com/SpecularL2/specular/blob/c54213cfb14aca9d44e839341f672dd978834f68/src/IRollup.sol)


## Functions
### assertions


```solidity
function assertions() external view returns (AssertionMap);
```

### isStaked


```solidity
function isStaked(address addr) external view returns (bool);
```
**Parameters**

|Name|Type|Description|
|----|----|-----------|
|`addr`|`address`|User address.|

**Returns**

|Name|Type|Description|
|----|----|-----------|
|`<none>`|`bool`|True if address is staked, else False.|


### currentRequiredStake


```solidity
function currentRequiredStake() external view returns (uint256);
```
**Returns**

|Name|Type|Description|
|----|----|-----------|
|`<none>`|`uint256`|The current required stake amount.|


### confirmedInboxSize


```solidity
function confirmedInboxSize() external view returns (uint256);
```
**Returns**

|Name|Type|Description|
|----|----|-----------|
|`<none>`|`uint256`|confirmedInboxSize size of inbox confirmed|


### stake

Deposits stake on staker's current assertion (or the last confirmed assertion if not currently staked).

currently use Ether to stake; stakeAmount Token amount to deposit. Must be > than defined threshold if this is a new stake.


```solidity
function stake() external payable;
```

### unstake

Withdraws stakeAmount from staker's stake by if assertion it is staked on is confirmed.


```solidity
function unstake(uint256 stakeAmount) external;
```
**Parameters**

|Name|Type|Description|
|----|----|-----------|
|`stakeAmount`|`uint256`|Token amount to withdraw. Must be <= sender's current stake minus the current required stake.|


### removeStake

Removes stakerAddress from the set of stakers and withdraws the full stake amount to stakerAddress.
This can be called by anyone since it is currently necessary to keep the chain progressing.


```solidity
function removeStake(address stakerAddress) external;
```
**Parameters**

|Name|Type|Description|
|----|----|-----------|
|`stakerAddress`|`address`|Address of staker for which to unstake.|


### advanceStake

Advances msg.sender's existing stake to assertionID.


```solidity
function advanceStake(uint256 assertionID) external;
```
**Parameters**

|Name|Type|Description|
|----|----|-----------|
|`assertionID`|`uint256`|ID of assertion to advance stake to. Currently this must be a child of the current assertion. TODO: generalize to arbitrary descendants.|


### withdraw

Withdraws all of msg.sender's withdrawable funds.


```solidity
function withdraw() external;
```

### createAssertion

Creates a new DA representing the rollup state after executing a block of transactions (sequenced in SequencerInbox).
Block is represented by all transactions in range [prevInboxSize, inboxSize]. The latest staked DA of the sender
is considered to be the predecessor. Moves sender stake onto the new DA.
The new DA stores the hash of the parameters: H(l2GasUsed || vmHash)


```solidity
function createAssertion(
    bytes32 vmHash,
    uint256 inboxSize,
    uint256 l2GasUsed,
    bytes32 prevVMHash,
    uint256 prevL2GasUsed
) external;
```
**Parameters**

|Name|Type|Description|
|----|----|-----------|
|`vmHash`|`bytes32`|New VM hash.|
|`inboxSize`|`uint256`|Size of inbox corresponding to assertion (number of transactions).|
|`l2GasUsed`|`uint256`|Total L2 gas used as of the end of this assertion's last transaction.|
|`prevVMHash`|`bytes32`|Predecessor assertion VM hash (required because it does not get stored in the assertion).|
|`prevL2GasUsed`|`uint256`|Predecessor assertion L2 gas used (required because it does not get stored in the assertion).|


### challengeAssertion

Initiates a dispute between a defender and challenger on an unconfirmed DA.


```solidity
function challengeAssertion(address[2] calldata players, uint256[2] calldata assertionIDs) external returns (address);
```
**Parameters**

|Name|Type|Description|
|----|----|-----------|
|`players`|`address[2]`|Defender (first) and challenger (second) addresses. Must be staked on DAs on different branches.|
|`assertionIDs`|`uint256[2]`|Assertion IDs of the players engaged in the challenge. The first ID should be the earlier-created and is the one being challenged.|

**Returns**

|Name|Type|Description|
|----|----|-----------|
|`<none>`|`address`|Newly created challenge contract address.|


### confirmFirstUnresolvedAssertion

Confirms first unresolved assertion. Assertion is confirmed if and only if:
(1) there is at least one staker, and
(2) challenge period has passed, and
(3) predecessor has been confirmed, and
(4) all stakers are staked on the assertion.


```solidity
function confirmFirstUnresolvedAssertion() external;
```

### rejectFirstUnresolvedAssertion

Rejects first unresolved assertion. Assertion is rejected if and only if:
(1) all of the following are true:
(a) challenge period has passed, and
(b) at least one staker exists, and
(c) no staker remains staked on the assertion (all have been destroyed).
OR
(2) predecessor has been rejected


```solidity
function rejectFirstUnresolvedAssertion(address stakerAddress) external;
```
**Parameters**

|Name|Type|Description|
|----|----|-----------|
|`stakerAddress`|`address`|Address of a staker staked on a different branch to the first unresolved assertion. If the first unresolved assertion's parent is confirmed, this parameter is used to establish that a staker exists on a different branch of the assertion chain. This parameter is ignored when the parent of the first unresolved assertion is not the last confirmed assertion.|


### completeChallenge

Completes ongoing challenge. Callback, called by a challenge contract.


```solidity
function completeChallenge(address winner, address loser) external;
```
**Parameters**

|Name|Type|Description|
|----|----|-----------|
|`winner`|`address`|Address of winning staker.|
|`loser`|`address`|Address of losing staker.|


## Events
### AssertionCreated

```solidity
event AssertionCreated(uint256 assertionID, address asserterAddr, bytes32 vmHash, uint256 inboxSize, uint256 l2GasUsed);
```

### AssertionChallenged

```solidity
event AssertionChallenged(uint256 assertionID, address challengeAddr);
```

### AssertionConfirmed

```solidity
event AssertionConfirmed(uint256 assertionID);
```

### AssertionRejected

```solidity
event AssertionRejected(uint256 assertionID);
```

### StakerStaked

```solidity
event StakerStaked(address stakerAddr, uint256 assertionID);
```

## Errors
### NotStaked
*Thrown when address that have not staked any token calls a only-staked function*


```solidity
error NotStaked();
```

### InsufficientStake
*Thrown when the function is called with Insufficient Stake*


```solidity
error InsufficientStake();
```

### StakedOnUnconfirmedAssertion
*Thrown when the caller is staked on unconfirmed assertion.*


```solidity
error StakedOnUnconfirmedAssertion();
```

### TransferFailed
*Thrown when transfer fails*


```solidity
error TransferFailed();
```

### AssertionOutOfRange
*Thrown when a staker tries to advance stake to invalid assertionId.*


```solidity
error AssertionOutOfRange();
```

### ParentAssertionUnstaked
*Thrown when a staker tries to advance stake to non-child assertion*


```solidity
error ParentAssertionUnstaked();
```

### MinimumAssertionPeriodNotPassed
*Thrown when a sender tries to create assertion before the minimum assertion time period*


```solidity
error MinimumAssertionPeriodNotPassed();
```

### MaxGasLimitExceeded
*Thrown when the L2 gas used by the assertion is more the max allowed limit.*


```solidity
error MaxGasLimitExceeded();
```

### PreviousStateHash
*Thrown when parent's statehash is not equal to the start state(or previous state)/*


```solidity
error PreviousStateHash();
```

### EmptyAssertion
*Thrown when a sender tries to create assertion without any tx.*


```solidity
error EmptyAssertion();
```

### InboxReadLimitExceeded
*Thrown when the requested assertion read past the end of current Inbox.*


```solidity
error InboxReadLimitExceeded();
```

### WrongOrder
*Thrown when the challenge assertion Id is not ordered or in range.*


```solidity
error WrongOrder();
```

### UnproposedAssertion
*Thrown when the challenger tries to challenge an unproposed assertion*


```solidity
error UnproposedAssertion();
```

### AssertionAlreadyResolved
*Thrown when the assertion is already resolved*


```solidity
error AssertionAlreadyResolved();
```

### NoUnresolvedAssertion
*Thrown when there is no unresolved assertion*


```solidity
error NoUnresolvedAssertion();
```

### ChallengePeriodPending
*Thrown when the challenge period has not passed*


```solidity
error ChallengePeriodPending();
```

### DifferentParent
*Thrown when the challenger and defender didn't attest to sibling assertions*


```solidity
error DifferentParent();
```

### InvalidParent
*Thrown when the assertion's parent is not the last confirmed assertion*


```solidity
error InvalidParent();
```

### NotInChallenge
*Thrown when the staker is not in a challenge*


```solidity
error NotInChallenge();
```

### InDifferentChallenge
*Thrown when the two stakers are in different challenge*


```solidity
error InDifferentChallenge(address staker1Challenge, address staker2Challenge);
```

### ChallengedStaker
*Thrown when the staker is currently in Challenge*


```solidity
error ChallengedStaker();
```

### NotAllStaked
*Thrown when all the stakers are not staked*


```solidity
error NotAllStaked();
```

### StakerStakedOnTarget
*Thrown staker's assertion is descendant of firstUnresolved assertion*


```solidity
error StakerStakedOnTarget();
```

### StakersPresent
*Thrown when there are staker's present on the assertion*


```solidity
error StakersPresent();
```

### NoStaker
*Thrown when there are zero stakers*


```solidity
error NoStaker();
```


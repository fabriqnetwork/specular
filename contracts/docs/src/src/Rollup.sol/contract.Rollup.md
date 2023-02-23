# Rollup
[Git Source](https://github.com/SpecularL2/specular/blob/559c78f8b09496c7f5c8f6e0b0262bee5e41d9a4/src/Rollup.sol)

**Inherits:**
[RollupBase](/src/Rollup.sol/contract.RollupBase.md)


## State Variables
### lastResolvedAssertionID

```solidity
uint256 public lastResolvedAssertionID;
```


### lastConfirmedAssertionID

```solidity
uint256 public lastConfirmedAssertionID;
```


### lastCreatedAssertionID

```solidity
uint256 public lastCreatedAssertionID;
```


### numStakers

```solidity
uint256 public numStakers;
```


### stakers

```solidity
mapping(address => Staker) public stakers;
```


### withdrawableFunds

```solidity
mapping(address => uint256) public withdrawableFunds;
```


### zombies

```solidity
Zombie[] public zombies;
```


## Functions
### stakedOnly


```solidity
modifier stakedOnly();
```

### initialize


```solidity
function initialize(
    address _vault,
    address _sequencerInbox,
    address _verifier,
    address _stakeToken,
    uint256 _confirmationPeriod,
    uint256 _challengePeriod,
    uint256 _minimumAssertionPeriod,
    uint256 _maxGasPerAssertion,
    uint256 _baseStakeAmount,
    bytes32 _initialVMhash
) public initializer;
```

### constructor


```solidity
constructor();
```

### _authorizeUpgrade


```solidity
function _authorizeUpgrade(address) internal override onlyOwner;
```

### isStaked


```solidity
function isStaked(address addr) public view override returns (bool);
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
function currentRequiredStake() public view override returns (uint256);
```
**Returns**

|Name|Type|Description|
|----|----|-----------|
|`<none>`|`uint256`|The current required stake amount.|


### confirmedInboxSize


```solidity
function confirmedInboxSize() public view override returns (uint256);
```
**Returns**

|Name|Type|Description|
|----|----|-----------|
|`<none>`|`uint256`|confirmedInboxSize size of inbox confirmed|


### stake

Deposits stake on staker's current assertion (or the last confirmed assertion if not currently staked).


```solidity
function stake() external payable override;
```

### unstake

Withdraws stakeAmount from staker's stake by if assertion it is staked on is confirmed.


```solidity
function unstake(uint256 stakeAmount) external override;
```
**Parameters**

|Name|Type|Description|
|----|----|-----------|
|`stakeAmount`|`uint256`|Token amount to withdraw. Must be <= sender's current stake minus the current required stake.|


### removeStake

Removes stakerAddress from the set of stakers and withdraws the full stake amount to stakerAddress.
This can be called by anyone since it is currently necessary to keep the chain progressing.


```solidity
function removeStake(address stakerAddress) external override;
```
**Parameters**

|Name|Type|Description|
|----|----|-----------|
|`stakerAddress`|`address`|Address of staker for which to unstake.|


### advanceStake

Advances msg.sender's existing stake to assertionID.


```solidity
function advanceStake(uint256 assertionID) external override stakedOnly;
```
**Parameters**

|Name|Type|Description|
|----|----|-----------|
|`assertionID`|`uint256`|ID of assertion to advance stake to. Currently this must be a child of the current assertion. TODO: generalize to arbitrary descendants.|


### withdraw

Withdraws all of msg.sender's withdrawable funds.


```solidity
function withdraw() external override;
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
) external override stakedOnly;
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


```solidity
function challengeAssertion(address[2] calldata players, uint256[2] calldata assertionIDs)
    external
    override
    returns (address);
```

### confirmFirstUnresolvedAssertion

Confirms first unresolved assertion. Assertion is confirmed if and only if:
(1) there is at least one staker, and
(2) challenge period has passed, and
(3) predecessor has been confirmed, and
(4) all stakers are staked on the assertion.


```solidity
function confirmFirstUnresolvedAssertion() external override;
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
function rejectFirstUnresolvedAssertion(address stakerAddress) external override;
```
**Parameters**

|Name|Type|Description|
|----|----|-----------|
|`stakerAddress`|`address`|Address of a staker staked on a different branch to the first unresolved assertion. If the first unresolved assertion's parent is confirmed, this parameter is used to establish that a staker exists on a different branch of the assertion chain. This parameter is ignored when the parent of the first unresolved assertion is not the last confirmed assertion.|


### completeChallenge

Completes ongoing challenge. Callback, called by a challenge contract.


```solidity
function completeChallenge(address winner, address loser) external override;
```
**Parameters**

|Name|Type|Description|
|----|----|-----------|
|`winner`|`address`|Address of winning staker.|
|`loser`|`address`|Address of losing staker.|


### stakeOnAssertion

Updates staker and assertion metadata.


```solidity
function stakeOnAssertion(address stakerAddress, uint256 assertionID) private;
```
**Parameters**

|Name|Type|Description|
|----|----|-----------|
|`stakerAddress`|`address`|Address of existing staker.|
|`assertionID`|`uint256`|ID of existing assertion to stake on.|


### deleteStaker

Deletes the staker from global state. Does not touch assertion staker state.


```solidity
function deleteStaker(address stakerAddress) private;
```
**Parameters**

|Name|Type|Description|
|----|----|-----------|
|`stakerAddress`|`address`|Address of the staker to delete|


### getChallenge

Checks to see whether the two stakers are in the same challenge


```solidity
function getChallenge(address staker1Address, address staker2Address) private view returns (address);
```
**Parameters**

|Name|Type|Description|
|----|----|-----------|
|`staker1Address`|`address`|Address of the first staker|
|`staker2Address`|`address`|Address of the second staker|

**Returns**

|Name|Type|Description|
|----|----|-----------|
|`<none>`|`address`|Address of the challenge that the two stakers are in|


### newAssertionDeadline


```solidity
function newAssertionDeadline() private view returns (uint256);
```

### countStakedZombies

Removes any zombies whose latest stake is earlier than the first unresolved assertion.

Counts the number of zombies staked on an assertion.

*Uses pop() instead of delete to prevent gaps, although order is not preserved*

*O(n), where n is # of zombies (but is expected to be small).
This function could be uncallable if there are too many zombies. However,
removeOldZombies() can be used to remove any zombies that exist so that this
will then be callable.*


```solidity
function countStakedZombies(uint256 assertionID) private view returns (uint256);
```
**Parameters**

|Name|Type|Description|
|----|----|-----------|
|`assertionID`|`uint256`|The assertion on which to count staked zombies|

**Returns**

|Name|Type|Description|
|----|----|-----------|
|`<none>`|`uint256`|The number of zombies staked on the assertion|


### requireStaked


```solidity
function requireStaked(address stakerAddress) private view;
```

### requireUnchallengedStaker


```solidity
function requireUnchallengedStaker(address stakerAddress) private view;
```


# SequencerInbox
[Git Source](https://github.com/SpecularL2/specular/blob/c54213cfb14aca9d44e839341f672dd978834f68/src/SequencerInbox.sol)

**Inherits:**
[ISequencerInbox](/src/ISequencerInbox.sol/contract.ISequencerInbox.md), Initializable, UUPSUpgradeable, OwnableUpgradeable


## State Variables
### inboxSize

```solidity
uint256 private inboxSize;
```


### accumulators

```solidity
bytes32[] public accumulators;
```


### sequencerAddress

```solidity
address public sequencerAddress;
```


## Functions
### initialize


```solidity
function initialize(address _sequencerAddress) public initializer;
```

### constructor


```solidity
constructor();
```

### _authorizeUpgrade


```solidity
function _authorizeUpgrade(address) internal override onlyOwner;
```

### getInboxSize


```solidity
function getInboxSize() external view override returns (uint256);
```

### appendTxBatch


```solidity
function appendTxBatch(uint256[] calldata contexts, uint256[] calldata txLengths, bytes calldata txBatch)
    external
    override;
```

### verifyTxInclusion


```solidity
function verifyTxInclusion(bytes memory proof) external view override;
```


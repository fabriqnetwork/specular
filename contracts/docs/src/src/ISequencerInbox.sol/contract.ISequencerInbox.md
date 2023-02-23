# ISequencerInbox
[Git Source](https://github.com/SpecularL2/specular/blob/559c78f8b09496c7f5c8f6e0b0262bee5e41d9a4/src/ISequencerInbox.sol)


## Functions
### getInboxSize

Gets inbox size (number of messages).


```solidity
function getInboxSize() external view returns (uint256);
```

### appendTxBatch

Appends a batch of transactions (stored in calldata) and emits a TxBatchAppended event.


```solidity
function appendTxBatch(uint256[] calldata contexts, uint256[] calldata txLengths, bytes calldata txBatch) external;
```
**Parameters**

|Name|Type|Description|
|----|----|-----------|
|`contexts`|`uint256[]`|Array of contexts, where each context is represented by a uint256 3-tuple: (numTxs, l2BlockNumber, l2Timestamp). Each context corresponds to a single "L2 block".|
|`txLengths`|`uint256[]`|Array of lengths of each encoded tx in txBatch.|
|`txBatch`|`bytes`|Batch of RLP-encoded transactions.|


### verifyTxInclusion

Verifies that a transaction exists in a batch, at the expected offset.


```solidity
function verifyTxInclusion(bytes memory proof) external view;
```
**Parameters**

|Name|Type|Description|
|----|----|-----------|
|`proof`|`bytes`|Proof of inclusion of transaction, in the form: proof := txInfo || batchInfo || {foreach tx in batch: (prefixHash || txDataHash), ...} where, txInfo := (sender || l2BlockNumber || l2Timestamp || txDataLength || txData) batchInfo := (batchNum || numTxsBefore || numTxsAfterInBatch || accBefore) TODO: modify based on OSP format.|


## Events
### TxBatchAppended

```solidity
event TxBatchAppended(uint256 batchNumber, uint256 startTxNumber, uint256 endTxNumber);
```

## Errors
### IncorrectAccOrBatch
*Thrown when the given tx inlcusion proof has incorrect accumulator or batch no.*


```solidity
error IncorrectAccOrBatch();
```

### EmptyBatch
*Thrown when sequencer tries to append an empty batch*


```solidity
error EmptyBatch();
```

### TxBatchDataOverflow
*Thrown when overflow occurs reading txBatch (likely due to malformed txLengths)*


```solidity
error TxBatchDataOverflow();
```


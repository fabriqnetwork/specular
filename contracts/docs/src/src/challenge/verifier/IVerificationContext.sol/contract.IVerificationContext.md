# IVerificationContext
[Git Source](https://github.com/SpecularL2/specular/blob/559c78f8b09496c7f5c8f6e0b0262bee5e41d9a4/src/challenge/verifier/IVerificationContext.sol)


## Functions
### getBlockHash


```solidity
function getBlockHash(uint8 number) external view returns (bytes32);
```

### getCoinbase


```solidity
function getCoinbase() external view returns (address);
```

### getTimestamp


```solidity
function getTimestamp() external view returns (uint256);
```

### getBlockNumber


```solidity
function getBlockNumber() external view returns (uint256);
```

### getDifficulty


```solidity
function getDifficulty() external view returns (uint256);
```

### getGasLimit


```solidity
function getGasLimit() external view returns (uint64);
```

### getChainID


```solidity
function getChainID() external view returns (uint256);
```

### getBaseFee


```solidity
function getBaseFee() external view returns (uint256);
```

### getStateRoot


```solidity
function getStateRoot() external view returns (bytes32);
```

### getEndStateRoot


```solidity
function getEndStateRoot() external view returns (bytes32);
```

### getOrigin


```solidity
function getOrigin() external view returns (address);
```

### getRecipient


```solidity
function getRecipient() external view returns (address);
```

### getTxnType


```solidity
function getTxnType() external view returns (TxnType);
```

### getValue


```solidity
function getValue() external view returns (uint256);
```

### getGas


```solidity
function getGas() external view returns (uint256);
```

### getGasPrice


```solidity
function getGasPrice() external view returns (uint256);
```

### getInput


```solidity
function getInput() external view returns (bytes memory);
```

### getInputSize


```solidity
function getInputSize() external view returns (uint64);
```

### getInputRoot


```solidity
function getInputRoot() external view returns (bytes32);
```

### getCodeMerkleFromInput


```solidity
function getCodeMerkleFromInput() external view returns (bytes32);
```

## Enums
### TxnType

```solidity
enum TxnType {
    TRANSFER,
    CREATE,
    CALL
}
```

